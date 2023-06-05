package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type Room struct {
	Title        string `json:"title"`
	Nickname     string `json:"nickname"`
	MaxUsers     int    `json:"maxUsers"`
	CurrentUsers int    `json:"currentUsers"`
}

var rooms []Room

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	server := socketio.NewServer(nil)

	// Socket.IO 이벤트 핸들러 등록
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("Socket.IO client connected:", s.ID())
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Socket.IO client disconnected:", s.ID())
	})

	// 방 추가 API 핸들러 등록
	router.POST("/rooms", func(c *gin.Context) {
		var newRoom Room
		if err := c.ShouldBindJSON(&newRoom); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 방 추가 로직
		rooms = append(rooms, newRoom)

		// 실시간으로 rooms 정보 업데이트
		server.BroadcastToRoom("/", "chat", "rooms", rooms)

		c.JSON(http.StatusOK, gin.H{"message": "Room created successfully"})
	})

	// Socket.IO 엔드포인트 등록
	server.OnEvent("/", "chat", func(s socketio.Conn, msg string) {
		log.Println("Received chat message:", msg)
	})

	go server.Serve()
	defer server.Close()

	// 프론트엔드에서 Socket.IO 클라이언트 스크립트 로드
	router.Static("/socket.io-client/", "./node_modules/socket.io-client/dist/")
	router.StaticFile("/", "index.html")

	router.NoRoute(func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	router.Run(":8080")
}
