package _interface

import (
	"database/sql"

	"github.com/saur4ig/file-storage/internal/models"
)

// FolderRepository - base functions to work with folders in postgres db
type FolderRepository interface {
	CreateFolder(userID int, name string, parentID int64) (int64, error)
	GetFolderByID(id int64) (*models.Folder, error)
	GetFoldersInfo(folderID int64) ([]models.FolderSize, error)
	// GetAllParentFolders returns all parent, and parent of parent folders
	GetAllParentFolders(folderID int64) ([]models.FolderSizeSimplified, error)
	DeleteFolder(tx *sql.Tx, id int64) error
	MoveFolder(tx *sql.Tx, folderID, newFolderID int64) error
	// UpdateFolderSize used to update the size only for this folder with new size
	UpdateFolderSize(id int64, newSize int64) error
	// IncreaseFolderSize used to add the size for this and all parent folders
	IncreaseFolderSize(tx *sql.Tx, id int64, size int64) error
	// DecreaseFolderSize used to reduce the size for this and all parent folders
	DecreaseFolderSize(tx *sql.Tx, id int64, size int64) error
	UpdateMultipleFoldersSize(tx *sql.Tx, folders []models.FolderSizeSimplified) error
}
