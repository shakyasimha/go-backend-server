// models/todo.go
package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id"` // uint to match gorm.Model.ID
}

func NewTodo() *Todo {
	return &Todo{}
}

func (t *Todo) ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sqlite3.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	if err := db.AutoMigrate(&Todo{}); err != nil {
		panic("Failed to automigrate Todo table: " + err.Error())
	}

	return db
}
