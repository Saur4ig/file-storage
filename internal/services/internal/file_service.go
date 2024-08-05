package internal

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	_interface "github.com/saur4ig/file-storage/internal/database/interface"
	"github.com/saur4ig/file-storage/internal/models"
	sinterface "github.com/saur4ig/file-storage/internal/services/interface"
)

type fileService struct {
	fileRepo   _interface.FileRepository
	folderRepo _interface.FolderRepository
	db         *sql.DB
}

func NewFileService(fileRepo _interface.FileRepository, folderRepo _interface.FolderRepository, db *sql.DB) sinterface.FileService {
	return &fileService{fileRepo: fileRepo, folderRepo: folderRepo, db: db}
}

// GetFile returns a file from the database
func (s *fileService) GetFile(fileID int64) (*models.File, error) {
	return s.fileRepo.GetFileByID(fileID)
}

// UploadFile uploads a file to a folder, updates folder size if necessary
func (s *fileService) UploadFile(folderID int64, userID int, name, s3URL string, size int64, transactionID *int64) error {
	// start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// transaction rollback in case of error
	defer func() {
		err = handleTxEnd(tx, err)
	}()

	file := &models.File{
		FolderID:      folderID,
		UserID:        userID,
		Name:          name,
		S3URL:         s3URL,
		Size:          size,
		TransactionID: transactionID,
	}

	// create file in db
	if err = s.fileRepo.CreateFile(tx, file); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	// if single file was added - update the size of folder and parent folders
	if transactionID == nil {
		if err = s.folderRepo.IncreaseFolderSize(tx, folderID, size); err != nil {
			return fmt.Errorf("failed to increase folder size: %w", err)
		}
	}

	return nil
}

// DeleteFile deletes a file and updates the folder size
func (s *fileService) DeleteFile(id int64) error {
	// start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// transaction rollback in case of error
	defer func() {
		err = handleTxEnd(tx, err)
	}()

	// get file data
	file, err := s.fileRepo.GetFileByID(id)
	if err != nil {
		return fmt.Errorf("failed to get file by ID: %w", err)
	}

	// remove the file
	if err = s.fileRepo.DeleteFile(tx, id); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return s.folderRepo.DecreaseFolderSize(tx, file.FolderID, file.Size)
}

// MoveFile moves a file to a new folder and updates the size of both folders.
func (s *fileService) MoveFile(fileID, folderID, newFolderID int64) error {
	// start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// transaction rollback in case of error
	defer func() {
		err = handleTxEnd(tx, err)
	}()

	// get file with data
	file, err := s.fileRepo.GetFileByID(fileID)
	if err != nil {
		return fmt.Errorf("failed to get file by ID: %w", err)
	}

	// change file folder
	err = s.fileRepo.MoveFile(tx, fileID, newFolderID)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	// decrease the old folder size
	if err = s.folderRepo.DecreaseFolderSize(tx, folderID, file.Size); err != nil {
		return fmt.Errorf("failed to decrease old folder size: %w", err)
	}

	// increase the new folder size
	return s.folderRepo.IncreaseFolderSize(tx, newFolderID, file.Size)
}

// handleTxEnd handles the end of a transaction, committing if no error, rolling back otherwise
func handleTxEnd(tx *sql.Tx, err error) error {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Warn().Msgf("Transaction rollback error: %v", rbErr)
		}
		return err
	}
	if cmErr := tx.Commit(); cmErr != nil {
		log.Warn().Msgf("Transaction commit error: %v", cmErr)
		return cmErr
	}
	return nil
}
