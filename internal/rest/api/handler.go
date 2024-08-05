package api

import (
	_interface "github.com/saur4ig/file-storage/internal/database/interface"
	si "github.com/saur4ig/file-storage/internal/services/interface"
)

type Handler struct {
	rc                 _interface.FolderSizeCache
	folderService      si.FolderService
	fileService        si.FileService
	transactionService si.TransactionService
	s3                 si.FileStorage
}

func New(
	fs si.FolderService,
	fileS si.FileService,
	ts si.TransactionService,
	s3 si.FileStorage,
	rc _interface.FolderSizeCache,
) *Handler {
	return &Handler{
		folderService:      fs,
		fileService:        fileS,
		transactionService: ts,
		s3:                 s3,
		rc:                 rc,
	}
}
