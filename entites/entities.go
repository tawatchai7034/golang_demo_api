package entites

import "time"

type Todo struct {
	Title     string `json:"text" binding:"required"`
	ID        uint   `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ResponseModel struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"reult"`
}
