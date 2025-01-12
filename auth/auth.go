package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
	Guid uuid.UUID
}

type Login struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (t *AuthHandler) Accesstoken(signature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body Login
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = json.Unmarshal(jsonData, &body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"massage": err.Error(),
				"result":  nil,
			})
			return
		}
		result := t.db.Where("Name = ?", body.User).First(&User{})
		fmt.Println(*result)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
			"sub": "user_id",
			"exp": time.Now().Add(15 * time.Minute).Unix(),
		})
		key := os.Getenv("JWT_SECRET_KEY")
		ss, err := token.SignedString([]byte(key))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"massage": err.Error(),
				"result":  nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"massage": "Success",
			"result":  ss,
		})
	}
}
