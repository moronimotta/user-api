package server

import (
	"log/slog"
	"time"
	userHandlers "user-auth/user/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) initializeMiddlewares() {
	// Middleware to log requests
	s.app.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		slog.Info("Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration,
		)
	})

	// Middleware to handle CORS
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		MaxAge:       12 * time.Hour,
	}))
}

// Authorization Middleware
func (s *Server) AuthMiddleware(userHttpHandler *userHandlers.UserHttpHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.Request.Header.Get("Role")
		token := c.Request.Header.Get("Authorization")
		if token == "" || userRole != "admin" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		ok, err := userHttpHandler.Repo.CheckAuthorizationRequest(userRole, token)
		if err != nil {
			slog.Error("Failed to check authorization request", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		// If token is valid, proceed to the next handler
		c.Next()
	}
}
