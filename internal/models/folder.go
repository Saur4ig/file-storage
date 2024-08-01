package models

import (
	"time"
)

// Folder represents a folder.
type Folder struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	Name           string    `db:"name"`
	ParentFolderID *int      `db:"parent_folder_id"`
	Size           int64     `db:"size"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
