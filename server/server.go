package server

import (
	"log/slog"
	"time"
	"user-auth/user/entities"
	userHandlers "user-auth/user/handlers"
	userRepo "user-auth/user/repositories"
	userUsecases "user-auth/user/usecases"

	"user-auth/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	app *gin.Engine
	db  db.Database
}

func NewServer(db db.Database) *Server {
	return &Server{
		app: gin.Default(),
		db:  db,
	}
}
func (s *Server) Start() {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))

	s.inicializeUserHttpHandler()

	if err := s.app.Run(":8080"); err != nil {
		panic(err)
	}
}

func (s *Server) inicializeUserHttpHandler() {
	userPostgresRepository := userRepo.NewUserPostgresRepository(s.db)
	userUsecase := userUsecases.NewUserUsecase(userPostgresRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(*userUsecase)

	userRoutes := s.app.Group("v1/user")
	userRoutes.POST("", func(c *gin.Context) {
		var input entities.User
		// get body
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Error("Failed to bind JSON", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := userHttpHandler.Repo.CreateUser(&input); err != nil {
			slog.Error("Failed to create user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "User created successfully"})
	})

	userRoutes.GET("", func(c *gin.Context) {
		user, err := userHttpHandler.Repo.GetAllUsers()
		if err != nil {
			slog.Error("Failed to get users", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"users": user})
	})

	userRoutes.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		user, err := userHttpHandler.Repo.GetUserByID(id)
		if err != nil {
			slog.Error("Failed to get user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"user": user})
	})

	userRoutes.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var input entities.User
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Error("Failed to bind JSON", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		input.ID = id
		if err := userHttpHandler.Repo.UpdateUser(&input); err != nil {
			slog.Error("Failed to update user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "User updated successfully"})
	})

	userRoutes.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := userHttpHandler.Repo.DeleteUser(id); err != nil {
			slog.Error("Failed to delete user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "User deleted successfully"})
	})

	userRoutes.GET("/account/:email", func(c *gin.Context) {
		email := c.Param("email")
		user, err := userHttpHandler.Repo.GetUserByEmail(email)
		if err != nil {
			slog.Error("Failed to get user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"user": user})
	})

	userRoutes.POST("/login", func(c *gin.Context) {
		var input entities.User
		if err := c.ShouldBindJSON(&input); err != nil {
			slog.Error("Failed to bind JSON", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		user, err := userHttpHandler.Repo.GetUserByEmailAndPassword(input.Email, input.Password)
		if err != nil {
			slog.Error("Failed to get user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(200, gin.H{"user": user})
	})

}
