package auth

import (
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
	store Storer
}

type Storer interface {
	Login(*Login) error
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	TransactionID() string
	Audience() string
}

func NewAuthHandler(store Storer) *AuthHandler {
	return &AuthHandler{store: store}
}

func (t AuthHandler) Accesstoken(c Context) {
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
