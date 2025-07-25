package server

import (
	"log/slog"
	"time"
	"user-auth/utils"

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

func (s *Server) ValidateRSAKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		ok := utils.ValidateJWTToken(c.Request.Header.Get("Authorization"))
		if !ok {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()

	}
}

func (s *Server) AdminPrivilegeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.Request.Header.Get("Role")
		if role != "admin" {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}
