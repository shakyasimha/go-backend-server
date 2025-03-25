package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Constructor, sort of?
func NewTodo() *Todo {
	return &Todo{}
}

// Function for connecting to db
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
