package database

import (
	"github.com/redis/go-redis/v9"
	_interface "github.com/saur4ig/file-storage/internal/database/interface"
	"github.com/saur4ig/file-storage/internal/database/internal"
)

func NewRedisCache(client *redis.Client) _interface.FolderSizeCache {
	return internal.NewRedisCache(client)
}
