package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/tawatchai7034/todo/auth"
	"github.com/tawatchai7034/todo/todo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// database handler
type DatabaseHandler struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Name string
	Guid uuid.UUID
}

func main() {
	//get .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// connect database
	db, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_CONNECTION")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&todo.Todo{}, &User{})
	// db.Create(&User{Name: "Fluke", Guid: uuid.New()})

	//defind path route
	userHandler := DatabaseHandler{db: db}
	r := gin.Default()
	todoHandler := todo.NewTodoHandler(db)
	authHandler := auth.NewAuthHandler(db)

	//middleware
	protected := r.Group("", auth.Protect([]byte(os.Getenv("JWT_SECRET_KEY"))))

	//get token
	r.POST("/login", authHandler.Accesstoken(os.Getenv("JWT_SECRET_KEY")))

	// router path
	protected.GET("/user", userHandler.User)
	protected.POST("/todo", todoHandler.NewTask)

	// Graceful Shutdown ,if not use Graceful Shutdown. You can change r.run()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%s\n", err)
		}
	}()

	// after get signal shutdown server 5 second
	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+c again to force")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

func (h *DatabaseHandler) User(c *gin.Context) {
	var u User
	h.db.First(&u)
	c.JSON(200, u)
}
