package db

import (
	"SafeTransfer/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Database holds the database connection.
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection.
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Automatically migrate your schema, ideally in development environments only.
	if err := db.AutoMigrate(&model.File{}); err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	log.Println("Successfully connected to database")
	return &Database{db}, nil
}

// Close gracefully closes the database connection.
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
