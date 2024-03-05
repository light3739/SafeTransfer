package db

import (
	"SafeTransfer/pkg/model"
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

	log.Println("Successfully connected to database")
	return &Database{db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// SaveFileMetadata saves the file metadata to the database.
func (d *Database) SaveFileMetadata(fileMetadata model.File) error {
	result := d.Create(&fileMetadata)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetFileMetadataByCID retrieves file metadata by CID.
func (d *Database) GetFileMetadataByCID(cid string) (*model.File, error) {
	var fileMetadata model.File
	result := d.Where("cid = ?", cid).First(&fileMetadata)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fileMetadata, nil
}
