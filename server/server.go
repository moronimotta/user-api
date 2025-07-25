package server

import (
	"log/slog"
	"time"
	"user-auth/user/entities"
	userHandlers "user-auth/user/handlers"
	userRepo "user-auth/user/repositories"
	userUsecases "user-auth/user/usecases"
	"user-auth/utils"

	"user-auth/db"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	app   *gin.Engine
	db    db.Database
	redis *redis.Client
}

func NewServer(db db.Database, redisClient *redis.Client) *Server {
	utils.InitLogging()

	return &Server{
		app:   gin.Default(),
		db:    db,
		redis: redisClient,
	}
}
func (s *Server) Start() {

	s.initializeMiddlewares()

	s.inicializeUserHttpHandler()

	if err := s.app.Run(":8080"); err != nil {
		panic(err)
	}
}

func (s *Server) inicializeUserHttpHandler() {
	userPostgresRepository := userRepo.NewUserPostgresRepository(s.db)
	userUsecase := userUsecases.NewUserUsecase(userPostgresRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(*userUsecase)

	rabbitMqHandler := userHandlers.NewRabbitMqHandler(s.db, s.redis)

	userRoutes := s.app.Group("v1/user")

	userRoutes.POST("", func(c *gin.Context) {
		var input entities.User
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

		dataInput := map[string]string{
			"email":   input.Email,
			"name":    input.Name,
			"user_id": input.ID,
		}

		if err := rabbitMqHandler.PublishMessage("finance-api", "user.created", dataInput); err != nil {
			slog.Error("Failed to publish message", err)
		}

		slog.Info("User created successfully", "user", input.Email)
		c.JSON(200, gin.H{"message": "User created successfully"})
	})

	// Route with middleware
	// s.AuthMiddleware(userHttpHandler),
	userRoutes.GET("", func(c *gin.Context) {
		user, err := userHttpHandler.Repo.GetAllUsers()
		if err != nil {
			slog.Error("Failed to get users", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		slog.Info("Get all users", "users")
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
		slog.Info("Get user by ID", "user", user.ID)
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

		dataInput := map[string]string{}
		if input.Email != "" {
			dataInput["email"] = input.Email
		}
		if input.Name != "" {
			dataInput["name"] = input.Name
		}

		if err := userHttpHandler.Repo.UpdateUser(&input); err != nil {
			slog.Error("Failed to update user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		updatedUser, err := userHttpHandler.Repo.GetUserByID(id)
		if err != nil {
			slog.Error("Failed to get updated user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		dataInput["external_id"] = updatedUser.ExternalID

		if err := rabbitMqHandler.PublishMessage("finance-api", "user.updated", dataInput); err != nil {
			slog.Error("Failed to publish message", err)
		}

		slog.Info("User updated successfully", "user", updatedUser.ID)
		c.JSON(200, gin.H{"message": "User updated successfully"})
	})

	userRoutes.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := userHttpHandler.Repo.DeleteUser(id); err != nil {
			slog.Error("Failed to delete user", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		slog.Info("User deleted successfully", "user", id)
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
		slog.Info("Get user by email", "user", user.Email)
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

		// generate JWT token
		token, err := utils.GenerateJWTWithRole(user.ID, user.Email, user.Role)
		if err != nil {
			slog.Error("Failed to generate JWT token", err)
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		// Set the JWT token in a cookie
		c.SetCookie("token", token, 3600, "/", "", false, true)
		// set cache on redis
		if err := s.redis.Set(c, user.ID, token, 3600*time.Second).Err(); err != nil {
			slog.Error("Failed to set cache on Redis", err)
		}

		slog.Info("User logged in successfully", "user", user.Email)
		c.JSON(200, gin.H{"user": user})
	})

}
