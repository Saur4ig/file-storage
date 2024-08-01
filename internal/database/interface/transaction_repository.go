package _interface

import (
	"github.com/saur4ig/file-storage/internal/models"
)

type TransactionRepository interface {
	CreateTransaction(tx *models.UploadTransaction) error
	GetTransactionByID(id int64) (*models.UploadTransaction, error)
	UpdateTransactionStatus(id int64, status string) error
}
