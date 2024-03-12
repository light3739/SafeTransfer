// FileRepository.go

package repository

import (
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/model"
	"gorm.io/gorm"
)

// FileRepository defines the interface for operations on the file entity.
type FileRepository interface {
	SaveFileMetadata(fileMetadata *model.File) error
	GetFileMetadataByCID(cid string) (*model.File, error)
}

// FileRepositoryImpl is the concrete implementation of FileRepository.
type FileRepositoryImpl struct {
	DB *gorm.DB
}

// NewFileRepository creates a new instance of FileRepositoryImpl.
func NewFileRepository(db *db.Database) FileRepository {
	return &FileRepositoryImpl{DB: db.DB}
}

// SaveFileMetadata saves the metadata of a file to the database.
func (repo *FileRepositoryImpl) SaveFileMetadata(fileMetadata *model.File) error {
	result := repo.DB.Create(fileMetadata)
	return result.Error
}

// GetFileMetadataByCID retrieves file metadata by CID.
func (repo *FileRepositoryImpl) GetFileMetadataByCID(cid string) (*model.File, error) {
	var fileMetadata model.File
	result := repo.DB.Where("cid = ?", cid).First(&fileMetadata)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fileMetadata, nil
}
