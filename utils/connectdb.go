package utils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ConnectDB connects to a SQLite database
func ConnectDB(dbName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	return db
}
