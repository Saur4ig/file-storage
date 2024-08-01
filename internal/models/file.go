package models

import (
	"time"
)

// File represents a file stored in the system.
type File struct {
	ID            int       `db:"id"`
	FolderID      int       `db:"folder_id"`
	UserID        int       `db:"user_id"`
	Name          string    `db:"name"`
	S3URL         string    `db:"s3_url"`
	Size          int64     `db:"size"`
	TransactionID *int      `db:"transaction_id"`
	CreatedAt     time.Time `db:"created_at"`
}
