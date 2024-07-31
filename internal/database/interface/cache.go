package _interface

import (
	"context"
)

// basic representation of the folders size cache
type FolderSizeCache interface {
	GetFolderSize(ctx context.Context, folderID int64) (int64, error)
	SetOrUpdateFolderSize(ctx context.Context, folderID int64, size int64) error
}
