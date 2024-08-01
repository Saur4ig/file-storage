package _interface

import (
	"github.com/saur4ig/file-storage/internal/models"
)

type FolderRepository interface {
	CreateFolder(userID int64, name string, parentID int64) error
	GetFolderByID(id int64) (*models.Folder, error)
	GetFoldersByParentID(parentFolderID int64) ([]*models.Folder, error)
	DeleteFolder(id int64) error
	UpdateFolderSize(id int64, newSize int64) error
}
