package services

import (
	"database/sql"

	rinterface "github.com/saur4ig/file-storage/internal/database/interface"
	_interface "github.com/saur4ig/file-storage/internal/services/interface"
	"github.com/saur4ig/file-storage/internal/services/internal"
)

func NewS3Service() _interface.FileStorage {
	return internal.NewS3Service()
}

func NewTransactionService(tr rinterface.TransactionRepository) _interface.TransactionService {
	return internal.NewTransactionService(tr)
}

func NewFileService(folderRepo rinterface.FolderRepository, fileRepo rinterface.FileRepository, db *sql.DB) _interface.FileService {
	return internal.NewFileService(fileRepo, folderRepo, db)
}

func NewFolderService(folderRepo rinterface.FolderRepository, fileRepo rinterface.FileRepository, db *sql.DB) _interface.FolderService {
	return internal.NewFolderService(folderRepo, fileRepo, db)
}
