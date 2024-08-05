package internal

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	_interface "github.com/saur4ig/file-storage/internal/database/interface"
)

type redisCache struct {
	client *redis.Client
}

type folderRepository struct {
	db *sql.DB
}

type transactionRepository struct {
	db *sql.DB
}

type fileRepository struct {
	db *sql.DB
}

func NewRedisCache(client *redis.Client) _interface.FolderSizeCache {
	return &redisCache{client: client}
}

func NewFolderRepository(db *sql.DB) _interface.FolderRepository {
	return &folderRepository{db: db}
}

func NewTransactionRepository(db *sql.DB) _interface.TransactionRepository {
	return &transactionRepository{db: db}
}

func NewFileRepository(db *sql.DB) _interface.FileRepository {
	return &fileRepository{db: db}
}
