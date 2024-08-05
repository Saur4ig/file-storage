package _interface

import (
	"github.com/saur4ig/file-storage/internal/models"
)

type FolderService interface {
	CreateFolder(userID int, name string, parentFolderID int64) (int64, error)
	DeleteFolder(id int64) error
	MoveFolder(folderID, newFolderID int64) error
	UpdateFolderSize(id int64, size int64) error
	GetFolderInfo(id int64) ([]models.FolderSize, error)
	GetAllParentFolders(folderID int64) ([]models.FolderSizeSimplified, error)
	UpdateMultipleFoldersSize(folders []models.FolderSizeSimplified) error
}
