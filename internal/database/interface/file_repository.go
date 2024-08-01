package _interface

import (
	"github.com/saur4ig/file-storage/internal/models"
)

type FileRepository interface {
	CreateFile(file *models.File) error
	GetFileByID(id int64) (*models.File, error)
	DeleteFile(id int64) error
}
