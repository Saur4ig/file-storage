package _interface

import (
	"context"

	"github.com/saur4ig/file-storage/internal/models"
)

// FolderSizeCache - basic representation of the folders size cache
type FolderSizeCache interface {
	GetFolderSize(ctx context.Context, folderID int64) (int64, error)
	SetOrUpdateFolderSize(ctx context.Context, folderID int64, size int64) error
	GetMultipleFolders(ctx context.Context, keys []string) ([]models.FolderSizeSimplified, error)
}
