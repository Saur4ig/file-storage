package internal

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/redis/go-redis/v9"
	"github.com/saur4ig/file-storage/internal/models"
)

const folderKeyPrefix = "folder_id:"

// GetFolderSize retrieves the size of a folder from Redis
func (rc *redisCache) GetFolderSize(ctx context.Context, folderID int64) (int64, error) {
	sizeStr, err := rc.client.Get(ctx, fmt.Sprintf("%s%d", folderKeyPrefix, folderID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil // Size not found in Redis
		}
		return 0, fmt.Errorf("error retrieving folder size: %w", err)
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing folder size: %w", err)
	}

	return size, nil
}

// SetOrUpdateFolderSize if the folder already in redis set, if not - update the size of a folder in Redis
func (rc *redisCache) SetOrUpdateFolderSize(ctx context.Context, folderID int64, size int64) error {
	key := fmt.Sprintf("%s%d", folderKeyPrefix, folderID)

	// Check if the folder size exists already
	sizeStr, err := rc.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return rc.setSize(ctx, key, size)
	}
	if err != nil {
		return fmt.Errorf("error checking folder size: %w", err)
	}

	return rc.updateSize(ctx, key, sizeStr, size)
}

// GetMultipleFolders returns all folders by provided keys
func (rc *redisCache) GetMultipleFolders(ctx context.Context, keys []string) ([]models.FolderSizeSimplified, error) {
	values, err := rc.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get folders from Redis: %w", err)
	}

	var folders []models.FolderSizeSimplified
	for i, val := range values {
		if val == nil {
			continue // Skip keys that do not exist
		}

		currentSize, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			log.Info().Msgf("error parsing folder size for key %s: %v", keys[i], err)
			continue
		}

		folderID, err := extractFolderID(keys[i])
		if err != nil {
			log.Info().Msgf("error extracting folderID for key %s: %v", keys[i], err)
			continue
		}

		folders = append(folders, models.FolderSizeSimplified{
			ID:   folderID,
			Size: currentSize,
		})
	}

	return folders, nil
}

// sets the folder size if it is not already set
func (rc *redisCache) setSize(ctx context.Context, key string, size int64) error {
	return rc.client.Set(ctx, key, size, 0).Err()
}

// updates the folder size if it already exists
func (rc *redisCache) updateSize(ctx context.Context, key string, sizeStr string, size int64) error {
	currentSize, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing current folder size: %w", err)
	}

	newSize := currentSize + size
	return rc.client.Set(ctx, key, newSize, 0).Err()
}

// extracts the folder ID from a Redis key formatted as "folder_id:<folderID>"
func extractFolderID(key string) (int64, error) {
	if !strings.HasPrefix(key, folderKeyPrefix) {
		return 0, fmt.Errorf("invalid key format: %s", key)
	}

	idStr := strings.TrimPrefix(key, folderKeyPrefix)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse folder ID from key: %w", err)
	}

	return id, nil
}
