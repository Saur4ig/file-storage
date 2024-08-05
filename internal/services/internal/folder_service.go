package internal

import (
	"database/sql"
	"fmt"

	rinterface "github.com/saur4ig/file-storage/internal/database/interface"
	"github.com/saur4ig/file-storage/internal/models"
	_interface "github.com/saur4ig/file-storage/internal/services/interface"
)

type folderService struct {
	fileRepo   rinterface.FileRepository
	folderRepo rinterface.FolderRepository
	db         *sql.DB
}

// NewFolderService creates a new FolderService
func NewFolderService(folderRepo rinterface.FolderRepository, fileRepo rinterface.FileRepository, db *sql.DB) _interface.FolderService {
	return &folderService{folderRepo: folderRepo, fileRepo: fileRepo, db: db}
}

// CreateFolder creates a new folder and returns its id
func (s *folderService) CreateFolder(userID int, name string, parentFolderID int64) (int64, error) {
	newFolderID, err := s.folderRepo.CreateFolder(userID, name, parentFolderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create folder: %w", err)
	}
	return newFolderID, nil
}

// MoveFolder moves a folder to a new parent folder and updates folder sizes accordingly
func (s *folderService) MoveFolder(folderID, newFolderID int64) error {
	folder, err := s.folderRepo.GetFolderByID(folderID)
	if err != nil {
		return fmt.Errorf("failed to get folder by ID: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		err = handleTxEnd(tx, err)
	}()

	oldFolderID := *folder.ParentFolderID

	// move the folder
	if err = s.folderRepo.MoveFolder(tx, folderID, newFolderID); err != nil {
		return fmt.Errorf("failed to move folder: %w", err)
	}

	// update sizes of the old and new parent folders
	if err = s.folderRepo.DecreaseFolderSize(tx, oldFolderID, folder.Size); err != nil {
		return fmt.Errorf("failed to decrease old folder size: %w", err)
	}

	if err = s.folderRepo.IncreaseFolderSize(tx, newFolderID, folder.Size); err != nil {
		return fmt.Errorf("failed to increase new folder size: %w", err)
	}

	return nil
}

// DeleteFolder deletes a folder and updates the parent folder size
func (s *folderService) DeleteFolder(id int64) error {
	folder, err := s.folderRepo.GetFolderByID(id)
	if err != nil {
		return fmt.Errorf("failed to get folder by ID: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		err = handleTxEnd(tx, err)
	}()

	if err = s.folderRepo.DeleteFolder(tx, id); err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	if err = s.folderRepo.DecreaseFolderSize(tx, *folder.ParentFolderID, folder.Size); err != nil {
		return fmt.Errorf("failed to decrease parent folder size: %w", err)
	}

	return nil
}

// UpdateFolderSize updates the size of a specified folder
func (s *folderService) UpdateFolderSize(id, size int64) error {
	if err := s.folderRepo.UpdateFolderSize(id, size); err != nil {
		return fmt.Errorf("failed to update folder size: %w", err)
	}
	return nil
}

// GetAllParentFolders retrieves all parent folders up to the root
func (s *folderService) GetAllParentFolders(folderID int64) ([]models.FolderSizeSimplified, error) {
	folders, err := s.folderRepo.GetAllParentFolders(folderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all parent folders: %w", err)
	}
	return folders, nil
}

// GetFolderInfo retrieves detailed information about a folder
func (s *folderService) GetFolderInfo(id int64) ([]models.FolderSize, error) {
	info, err := s.folderRepo.GetFoldersInfo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder info: %w", err)
	}
	return info, nil
}

// UpdateMultipleFoldersSize updates the sizes of multiple folders within a transaction
func (s *folderService) UpdateMultipleFoldersSize(folders []models.FolderSizeSimplified) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		err = handleTxEnd(tx, err)
	}()

	if err = s.folderRepo.UpdateMultipleFoldersSize(tx, folders); err != nil {
		return fmt.Errorf("failed to update multiple folders' sizes: %w", err)
	}

	return nil
}
