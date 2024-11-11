// sockets/chat.go
package sockets

import (
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func ChatSocket() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	server.OnEvent("/", "join", func(s socketio.Conn, room string) {
		s.Join(room)
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		server.BroadcastToRoom("/", s.Context().(string), "message", msg)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	return server
}

func ServeSocket(c *gin.Context) {
	server := ChatSocket()
	server.ServeHTTP(c.Writer, c.Request)
}
