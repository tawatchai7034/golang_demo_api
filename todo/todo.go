package todo

import (
	"net/http"
	"time"

	"github.com/tawatchai7034/todo/entites"
)

type Todo struct {
	Title     string `json:"text" binding:"required"`
	ID        uint   `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storer interface {
	New(*Todo) error
}

func (Todo) TableName() string {
	return "todoList"
}

type TodoHandler struct {
	store Storer
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	TransactionID() string
	Audience() string
}

func NewTodoHandler(store Storer) *TodoHandler {
	return &TodoHandler{store: store}
}

func (t *TodoHandler) NewTask(c Context) {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, entites.ResponseModel{
			Status:  "error",
			Message: err.Error(),
			Result:  nil,
		})
		return
	}

	// create table in database and row data
	err := t.store.New(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, entites.ResponseModel{
			Status:  "error",
			Message: err.Error(),
			Result:  nil,
		})
		return
	}

	c.JSON(http.StatusCreated, entites.ResponseModel{
		Status:  "Success",
		Message: "Success",
		Result:  todo.ID,
	})
}
