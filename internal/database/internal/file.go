package internal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/saur4ig/file-storage/internal/models"
)

// CreateFile inserts a new file record into the database
func (r *fileRepository) CreateFile(tx *sql.Tx, file *models.File) error {
	query := `
		INSERT INTO files (folder_id, user_id, name, s3_url, size, transaction_id) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at
	`
	if err := tx.QueryRow(query, file.FolderID, file.UserID, file.Name, file.S3URL, file.Size, file.TransactionID).
		Scan(&file.ID, &file.CreatedAt); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	return nil
}

// GetFileByID retrieves a file from the database by its id
func (r *fileRepository) GetFileByID(id int64) (*models.File, error) {
	query := `
		SELECT id, folder_id, user_id, name, s3_url, size, transaction_id, created_at 
		FROM files 
		WHERE id = $1
	`
	file := &models.File{}
	err := r.db.QueryRow(query, id).
		Scan(&file.ID, &file.FolderID, &file.UserID, &file.Name, &file.S3URL, &file.Size, &file.TransactionID, &file.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("file not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve file by ID: %w", err)
	}
	return file, nil
}

// DeleteFile deletes a file record from the database by its id
func (r *fileRepository) DeleteFile(tx *sql.Tx, id int64) error {
	query := `DELETE FROM files WHERE id = $1`
	if _, err := tx.Exec(query, id); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// MoveFile updates the folder_id of a file to move it to new folder
func (r *fileRepository) MoveFile(tx *sql.Tx, fileID, newFolderID int64) error {
	query := `UPDATE files SET folder_id = $1 WHERE id = $2`
	if _, err := tx.Exec(query, newFolderID, fileID); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	return nil
}
