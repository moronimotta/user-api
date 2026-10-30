package server

import (
	"errors"
	"log/slog"
	"strings"
	"time"
	"user-auth/user/entities"
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
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Origin"},
		MaxAge:       12 * time.Hour,
	}))
}

func (s *Server) ValidateRSAKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := parseRequestClaims(c); err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()

	}
}

func (s *Server) AuthenticatedUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := parseRequestClaims(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		id, _ := claims["id"].(string)
		subject, _ := claims["sub"].(string)
		if id == "" {
			id = subject
		}
		if id == "" || (subject != "" && subject != id) {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || strings.TrimSpace(role) == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("auth_user_id", id)
		c.Set("auth_user_role", entities.CanonicalRole(role))
		c.Next()
	}
}

func (s *Server) AdminPrivilegeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAdminRequest(c) {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func isAdminRequest(c *gin.Context) bool {
	claims, err := parseRequestClaims(c)
	if err != nil {
		return false
	}
	id, _ := claims["id"].(string)
	subject, _ := claims["sub"].(string)
	if id == "" {
		id = subject
	}
	if id == "" || (subject != "" && subject != id) {
		return false
	}
	role, ok := claims["role"].(string)
	return ok && entities.CanonicalRole(role) == entities.RoleAdmin
}

func parseRequestClaims(c *gin.Context) (map[string]interface{}, error) {
	token := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
	}
	if token == "" {
		token, _ = c.Cookie("token")
	}
	if token == "" {
		return nil, errors.New("missing authentication token")
	}
	claims, err := utils.ParseJWTTokenClaims(token)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}(claims), nil
}

func canAccessUser(c *gin.Context, userID string) bool {
	authUserID, _ := c.Get("auth_user_id")
	role, _ := c.Get("auth_user_role")
	return authUserID == userID || role == entities.RoleAdmin
}
