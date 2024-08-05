package models

import (
	"time"
)

// Folder represents a folder
type Folder struct {
	ID             int64     `db:"id"`
	UserID         int       `db:"user_id"`
	Name           string    `db:"name"`
	ParentFolderID *int64    `db:"parent_folder_id"`
	Size           int64     `db:"size"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// FolderSize represents the size of a folder and its metadata
type FolderSize struct {
	Name string `db:"name"`
	ID   int64  `db:"id"`
	Size int64  `db:"size"`
}

// FolderSizeSimplified represents the size of a folder and only id
type FolderSizeSimplified struct {
	ID   int64 `db:"id"`
	Size int64 `db:"size"`
}
