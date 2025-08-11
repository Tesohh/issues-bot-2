package db

import (
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Connect(path string) (*gorm.DB, error) {
	slog.Info("Connecting to db at", "path", path)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate()

	Conn = db
	slog.Info("Connected to db")
	return db, nil
}
