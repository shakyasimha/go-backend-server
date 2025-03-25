// models/user.go
package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"` // Will be hashed
	Email    string `json:"email"`
}

type Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser() *User {
	return &User{}
}

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

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}
