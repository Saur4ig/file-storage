package internal

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// GetFolderSize retrieves the size of a folder from Redis
func (rc *redisCache) GetFolderSize(ctx context.Context, folderID int64) (int64, error) {
	sizeStr, err := rc.client.Get(ctx, fmt.Sprintf("folder_id:%d", folderID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// size not found in Redis
			return 0, nil
		} else {
			return 0, err
		}
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// SetOrUpdateFolderSize if the folder already in redis set, if not - update the size of a folder in Redis
func (rc *redisCache) SetOrUpdateFolderSize(ctx context.Context, folderID int64, size int64) error {
	key := fmt.Sprintf("folder_id:%d", folderID)

	// check if the folder exists already
	sizeStr, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		// there is no folder -> set a size
		if errors.Is(err, redis.Nil) {
			return rc.setSize(ctx, key, size)
		} else {
			return err
		}
	} else {
		return rc.updateSize(ctx, key, sizeStr, size)
	}
}

// set if not found
func (rc *redisCache) setSize(ctx context.Context, key string, size int64) error {
	return rc.client.Set(ctx, key, size, 0).Err()
}

// update if exists
func (rc *redisCache) updateSize(ctx context.Context, key string, sizeStr string, size int64) error {
	currentSize, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return err
	}

	newSize := currentSize + size
	return rc.client.Set(ctx, key, newSize, 0).Err()
}
