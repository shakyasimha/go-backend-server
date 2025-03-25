package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Pseudoconstructor
func NewUser() *User {
	return &User{}
}

// Function for connecting to db and migrating
func (u *User) ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sqlite3.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		panic("Failed to automigrate User table: " + err.Error())
	}

	return db
}
