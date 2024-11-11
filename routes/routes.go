// router/router.go
package router

import (
	"yak.rabai/controllers"
	// "yak.rabai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Trust all proxies (for testing), but ideally configure for production
	router.SetTrustedProxies([]string{})

	// Public routes
	api := router.Group("/api")
	{
		// Authentication routes
		api.GET("/auth/google", controllers.GoogleLogin)
		api.GET("/auth/google/callback", controllers.GoogleCallback)

		// Queue routes
		api.POST("/queue", controllers.QueueChatRoom)

		// Chat WebSocket connection
		api.GET("/ws", controllers.ConnectChat)

		//TODO: Add middleware for authentication
		// Chat endpoints with authentication
		chat := api.Group("/chat")
		{
			chat.POST("/message", controllers.HandleChatMessage)    // Send message
			chat.POST("/typing", controllers.HandleTypingIndicator) // Typing indicator
			chat.POST("/leave", controllers.LeaveRoomAndRate)       // Leave room and rate listener
		}
	}

	return router
}
