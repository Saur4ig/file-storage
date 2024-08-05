package internal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/saur4ig/file-storage/internal/models"
)

// CreateTransaction inserts a new upload transaction into the database and returns the created transaction ID
func (r *transactionRepository) CreateTransaction(userID int, folderID int64) (int64, error) {
	const pendingStatus = "pending"
	var transactionID int64
	query := `
		INSERT INTO upload_transactions (user_id, folder_id, status) 
		VALUES ($1, $2, $3) 
		RETURNING id
	`
	err := r.db.QueryRow(query, userID, folderID, pendingStatus).Scan(&transactionID)
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}
	return transactionID, nil
}

// GetTransactionByID retrieves an upload transaction by its id
func (r *transactionRepository) GetTransactionByID(id int64) (*models.UploadTransaction, error) {
	query := `
		SELECT id, user_id, folder_id, status, created_at, updated_at 
		FROM upload_transactions 
		WHERE id = $1
	`
	tx := &models.UploadTransaction{}
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.UserID,
		&tx.FolderID,
		&tx.Status,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No transaction found
		}
		return nil, fmt.Errorf("failed to retrieve transaction by ID: %w", err)
	}
	return tx, nil
}

// UpdateTransactionStatus updates the status of a transaction
func (r *transactionRepository) UpdateTransactionStatus(id int64, status string) error {
	query := `
		UPDATE upload_transactions 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2
	`
	if _, err := r.db.Exec(query, status, id); err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}
	return nil
}
