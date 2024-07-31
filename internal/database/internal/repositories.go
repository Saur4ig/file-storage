package internal

import (
	"github.com/redis/go-redis/v9"
	_interface "github.com/saur4ig/file-storage/internal/database/interface"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) _interface.FolderSizeCache {
	return &redisCache{client: client}
}
