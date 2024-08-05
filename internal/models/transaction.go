package models

import (
	"time"
)

// UploadTransaction represents a file upload transaction.
type UploadTransaction struct {
	ID        int64     `db:"id"`
	UserID    int       `db:"user_id"`
	FolderID  int64     `db:"folder_id"`
	TotalSize int64     `db:"total_size"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
