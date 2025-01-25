package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tawatchai7034/todo/entites"
	"gorm.io/gorm"
)

type Todo struct {
	entites.Todo
}

func (Todo) TableName() string {
	return "todoList"
}

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": err.Error(),
			"result":  nil,
		})
		return
	}

	// create table in database and row data
	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"massage": err.Error(),
			"result":  nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Success",
		"massage": "Success",
		"result":  todo.Model.ID,
	})
}
