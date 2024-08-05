package _interface

import (
	"github.com/saur4ig/file-storage/internal/models"
)

type TransactionService interface {
	CreateTransaction(userID int, folderID int64) (int64, error)
	GetTransactionByID(id int64) (*models.UploadTransaction, error)
	UpdateTransactionStatus(id int64, status string) error
}
