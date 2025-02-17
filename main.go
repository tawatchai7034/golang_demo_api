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

	"golang.org/x/time/rate"

	"github.com/gin-contrib/cors"
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

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	//get .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Liveness Probe and create file
	f, err := os.Create("D:/Learning/golang_demo_api/tmp/live.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("D:/Learning/golang_demo_api/tmp/live.txt")
	defer f.Close()

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
	// - No origin allowed by default
	// - GET, POST, PUT, PATCH, DELETE HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("HOST")}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}
	r.Use(cors.New((config)))

	todoHandler := todo.NewTodoHandler(db)
	authHandler := auth.NewAuthHandler(db)

	//middleware
	protected := r.Group("", auth.Protect([]byte(os.Getenv("JWT_SECRET_KEY"))))

	//path get token
	r.POST("/login", authHandler.Accesstoken(os.Getenv("JWT_SECRET_KEY")))

	// router path
	protected.GET("/user", userHandler.User)
	protected.POST("/todo", todoHandler.NewTask)

	//path get build commit
	r.GET("/buildLog", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	//path Readiness Probe
	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "Success",
			"massage": "Success",
			"result":  nil,
		})
	})

	//path limiter handler
	r.GET("/limitz", limitedHandler)

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

// Rate Limit handler
var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(200, gin.H{
		"status":  "Success",
		"massage": "Success",
		"result":  nil,
	})
}
