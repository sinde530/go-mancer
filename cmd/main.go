package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ChatRoom struct {
	Name     string
	MaxUsers int
	Users    []string
}

type ChatApp struct {
	Rooms     map[string]*ChatRoom
	RoomsLock sync.Mutex
}

func main() {
	app := &ChatApp{
		Rooms: make(map[string]*ChatRoom),
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // 클라이언트의 도메인 주소
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	router.Use(cors.New(config))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the chat app!")
	})

	router.GET("/rooms", func(c *gin.Context) {
		app.RoomsLock.Lock()
		defer app.RoomsLock.Unlock()

		roomNames := make([]string, 0, len(app.Rooms))
		for roomName := range app.Rooms {
			roomNames = append(roomNames, roomName)
		}

		c.JSON(http.StatusOK, roomNames)

	})

	router.GET("/join/:roomName", func(c *gin.Context) {
		roomName := c.Param("roomName")
		c.String(http.StatusOK, "Joined as: %s", roomName)

		// 방 참여 동작 구현
		// roomName := "My Room" // 예시로 고정된 방 이름 사용
		app.RoomsLock.Lock()
		defer app.RoomsLock.Unlock()

		if room, ok := app.Rooms[roomName]; ok {
			room.Users = append(room.Users, roomName)
			c.String(http.StatusOK, "Joined room: %s", roomName)
		} else {
			c.String(http.StatusBadRequest, "Room not found")
		}
	})

	router.POST("/create-room", func(c *gin.Context) {
		var json struct {
			Name     string `json:"name"`
			MaxUsers int    `json:"max_users"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.String(http.StatusBadRequest, "Invalid request")
			return
		}

		app.RoomsLock.Lock()
		defer app.RoomsLock.Unlock()

		if _, ok := app.Rooms[json.Name]; ok {
			c.String(http.StatusBadRequest, "Room already exists")
			return
		}

		room := &ChatRoom{
			Name:     json.Name,
			MaxUsers: json.MaxUsers,
			Users:    []string{},
		}

		app.Rooms[json.Name] = room

		c.String(http.StatusOK, "Room created: %s", json.Name)
	})

	router.POST("/delete-room/:roomName", func(c *gin.Context) {
		roomName := c.Param("roomName")

		app.RoomsLock.Lock()
		defer app.RoomsLock.Unlock()

		if _, ok := app.Rooms[roomName]; ok {
			delete(app.Rooms, roomName)
			c.String(http.StatusOK, "Room deleted: %s", roomName)
		} else {
			c.String(http.StatusBadRequest, "Room not found")
		}
	})

	log.Fatal(router.Run(":8080"))
}
