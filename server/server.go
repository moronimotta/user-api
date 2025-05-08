package server

import (
	"time"
	dummy "user-auth/user/entities"
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

	dummyUser := dummy.UserEntity

	userRoutes := s.app.Group("v1/user")
	userRoutes.POST("", func(c *gin.Context) {
		if err := userHttpHandler.Repo.CreateUser(&dummyUser); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "User created successfully"})
	})

	userRoutes.GET("", func(c *gin.Context) {
		user, err := userHttpHandler.Repo.GetAllUsers()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"users": user})
	})

}
