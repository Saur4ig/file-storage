package models

import (
	"time"
)

// File represents a file stored in the system.
type File struct {
	ID            int64     `db:"id"`
	FolderID      int64     `db:"folder_id"`
	UserID        int       `db:"user_id"`
	Name          string    `db:"name"`
	S3URL         string    `db:"s3_url"`
	Size          int64     `db:"size"`
	TransactionID *int64    `db:"transaction_id"`
	CreatedAt     time.Time `db:"created_at"`
}
