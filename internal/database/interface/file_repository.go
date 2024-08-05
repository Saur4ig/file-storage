package _interface

import (
	"database/sql"

	"github.com/saur4ig/file-storage/internal/models"
)

// FileRepository - base functions to work with stored in postgres db file data
type FileRepository interface {
	CreateFile(tx *sql.Tx, file *models.File) error
	GetFileByID(id int64) (*models.File, error)
	DeleteFile(tx *sql.Tx, id int64) error
	MoveFile(tx *sql.Tx, fileID, newFolderID int64) error
}
