package store

import (
	"github.com/tawatchai7034/todo/auth"
	"github.com/tawatchai7034/todo/todo"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) New(todo *todo.Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) Login(todo *auth.Login) error {
	return nil
}
