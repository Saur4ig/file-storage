package database

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	_interface "github.com/saur4ig/file-storage/internal/database/interface"
	"github.com/saur4ig/file-storage/internal/database/internal"
)

func NewRedisCache(client *redis.Client) _interface.FolderSizeCache {
	return internal.NewRedisCache(client)
}

func NewFolderRepository(db *sql.DB) _interface.FolderRepository {
	return internal.NewFolderRepository(db)
}

func NewTransactionRepository(db *sql.DB) _interface.TransactionRepository {
	return internal.NewTransactionRepository(db)
}

func NewFileRepository(db *sql.DB) _interface.FileRepository {
	return internal.NewFileRepository(db)
}
