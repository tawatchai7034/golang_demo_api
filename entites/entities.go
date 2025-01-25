package entites

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title string `json:"text"`
}
