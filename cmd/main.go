package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID       string
	Clients  map[*websocket.Conn]bool
	Messages []string
}

func (r *Room) Broadcast(message string) {
	r.Messages = append(r.Messages, message)
	for client := range r.Clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Failed to send message to client:", err)
			client.Close()
			delete(r.Clients, client)
		}
	}
}

// func (r *Room) Broadcast(message string) {
// 	r.Messages = append(r.Messages, message)

// 	userCount := len(r.Clients)

// 	for client := range r.Clients {
// 		err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("message:%s,userCount:%d", message, userCount)))
// 		if err != nil {
// 			log.Println("Failed to send message to client:", err)
// 			client.Close()
// 			delete(r.Clients, client)
// 		}
// 	}
// }

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	rooms := make(map[string]*Room)
	roomIDs := []string{"room1", "room2"}
	for _, id := range roomIDs {
		rooms[id] = &Room{
			ID:       id,
			Clients:  make(map[*websocket.Conn]bool),
			Messages: []string{},
		}
	}

	router.GET("/ws/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to upgrade connection:", err)
			return
		}

		room, exists := rooms[roomID]
		if !exists {
			log.Println("Invalid room:", roomID)
			conn.Close()
			return
		}

		if len(room.Clients) >= 2 {
			log.Println("Room is full")
			conn.Close()
			return
		}

		room.Clients[conn] = true

		for _, message := range room.Messages {
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Failed to send message to client:", err)
				conn.Close()
				delete(room.Clients, conn)
				return
			}
		}

		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("Failed to read message from client:", err)
					conn.Close()
					delete(room.Clients, conn)
					return
				}
				room.Broadcast(string(msg))
			}
		}()
	})

	router.GET("/api/rooms/:roomID/users", func(c *gin.Context) {
		roomID := c.Param("roomID")

		room, exists := rooms[roomID]
		if !exists {
			log.Println("Invalid room:", roomID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid room"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"count": len(room.Clients)})
	})

	router.Run(":8080")
}
