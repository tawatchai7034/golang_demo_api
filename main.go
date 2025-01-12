package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	db, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_CONNECTION")), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&User{})
	// db.Create(&User{Name: "Fluke"})
	//defind path route
	userHandler := DatabaseHandler{db: db}
	r := gin.Default()
	r.GET("/user", userHandler.User)
	r.Run()
}

func (h *DatabaseHandler) User(c *gin.Context) {
	var u User
	h.db.First(&u)
	c.JSON(200, u)
}
