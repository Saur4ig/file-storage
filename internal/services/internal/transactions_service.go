package internal

import (
	"errors"

	rinterface "github.com/saur4ig/file-storage/internal/database/interface"
	"github.com/saur4ig/file-storage/internal/models"
	_interface "github.com/saur4ig/file-storage/internal/services/interface"
)

type transactionService struct {
	transactionRepo rinterface.TransactionRepository
}

func NewTransactionService(transactionRepo rinterface.TransactionRepository) _interface.TransactionService {
	return &transactionService{transactionRepo: transactionRepo}
}

func (s *transactionService) CreateTransaction(userID int, folderID int64) (int64, error) {
	return s.transactionRepo.CreateTransaction(userID, folderID)
}

func (s *transactionService) GetTransactionByID(id int64) (*models.UploadTransaction, error) {
	return s.transactionRepo.GetTransactionByID(id)
}

func (s *transactionService) UpdateTransactionStatus(id int64, status string) error {
	if status != "pending" && status != "completed" && status != "failed" {
		return errors.New("invalid status")
	}
	return s.transactionRepo.UpdateTransactionStatus(id, status)
}
