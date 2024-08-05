package _interface

import (
	"github.com/saur4ig/file-storage/internal/models"
)

type FileService interface {
	GetFile(fileID int64) (*models.File, error)
	UploadFile(folderID int64, userID int, name, s3URL string, size int64, transactionID *int64) error
	MoveFile(fileID, folderID, newFolderID int64) error
	DeleteFile(id int64) error
}
