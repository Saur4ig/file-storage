package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/saur4ig/file-storage/internal/models"
)

// CreateFolder creates a folder and returns its id, if successful.
func (r *folderRepository) CreateFolder(userID int, name string, parentID int64) (int64, error) {
	query := `INSERT INTO folders (user_id, name, parent_folder_id) VALUES ($1, $2, $3) RETURNING id`
	var folderID int64
	err := r.db.QueryRow(query, userID, name, parentID).Scan(&folderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create folder: %w", err)
	}
	return folderID, nil
}

// GetFolderByID retrieves all data of a folder by its id.
func (r *folderRepository) GetFolderByID(id int64) (*models.Folder, error) {
	query := `SELECT id, user_id, name, parent_folder_id, size, created_at, updated_at FROM folders WHERE id = $1`
	folder := &models.Folder{}
	err := r.db.QueryRow(query, id).Scan(&folder.ID, &folder.UserID, &folder.Name, &folder.ParentFolderID, &folder.Size, &folder.CreatedAt, &folder.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("folder not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve folder by ID: %w", err)
	}
	return folder, nil
}

// GetFoldersInfo retrieves the folder and it all parent subfolders sizes
func (r *folderRepository) GetFoldersInfo(folderID int64) ([]models.FolderSize, error) {
	query := `
		SELECT id, name, size FROM folders WHERE id = $1
		UNION
		SELECT id, name, size FROM folders WHERE parent_folder_id = $1
	`
	rows, err := r.db.Query(query, folderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve folders info: %w", err)
	}
	defer rows.Close()

	var folders []models.FolderSize
	for rows.Next() {
		var folderWithSize models.FolderSize
		if err := rows.Scan(&folderWithSize.ID, &folderWithSize.Name, &folderWithSize.Size); err != nil {
			return nil, fmt.Errorf("failed to scan folder size: %w", err)
		}
		folders = append(folders, folderWithSize)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating folder rows: %w", err)
	}

	return folders, nil
}

// GetAllParentFolders retrieves all parent folders up to the root folder
func (r *folderRepository) GetAllParentFolders(folderID int64) ([]models.FolderSizeSimplified, error) {
	query := `
		WITH RECURSIVE parent_folders AS (
			SELECT id, size, parent_folder_id
			FROM folders
			WHERE id = $1
			UNION ALL
			SELECT f.id, f.size, f.parent_folder_id
			FROM folders f
			INNER JOIN parent_folders pf ON f.id = pf.parent_folder_id
		)
		SELECT id, size
		FROM parent_folders
		ORDER BY id
	`
	rows, err := r.db.Query(query, folderID)
	if err != nil {
		log.Printf("Error querying parent folders up to root: %v", err)
		return nil, fmt.Errorf("failed to retrieve parent folders: %w", err)
	}
	defer rows.Close()

	var folders []models.FolderSizeSimplified
	for rows.Next() {
		var folderWithSize models.FolderSizeSimplified
		if err := rows.Scan(&folderWithSize.ID, &folderWithSize.Size); err != nil {
			return nil, fmt.Errorf("failed to scan parent folder size: %w", err)
		}
		folders = append(folders, folderWithSize)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating parent folder rows: %w", err)
	}

	return folders, nil
}

// DeleteFolder removes folder and all contents inside
func (r *folderRepository) DeleteFolder(tx *sql.Tx, id int64) error {
	// recursively delete the folder data
	return r.recursiveDeleteFolder(tx, id)
}

// MoveFolder moves a folder to another folder, ensuring no cycles are created
func (r *folderRepository) MoveFolder(tx *sql.Tx, folderID, newFolderID int64) error {
	// check if moving folder creates a cycle
	if err := r.checkForCycle(tx, folderID, newFolderID); err != nil {
		return err
	}

	// update the folder's parent_folder_id
	_, err := tx.Exec(`UPDATE folders SET parent_folder_id = $1 WHERE id = $2`, newFolderID, folderID)
	if err != nil {
		return fmt.Errorf("failed to move folder: %w", err)
	}

	return nil
}

// UpdateFolderSize replaces actual size of the folder with new size only
// used for updating the size after multiple files upload
func (r *folderRepository) UpdateFolderSize(id, newSize int64) error {
	query := `UPDATE folders SET size = $1, updated_at = NOW() WHERE id = $2`
	if _, err := r.db.Exec(query, newSize, id); err != nil {
		return fmt.Errorf("failed to update folder size: %w", err)
	}
	return nil
}

// IncreaseFolderSize increases a folder size and the size of all parent folders
func (r *folderRepository) IncreaseFolderSize(tx *sql.Tx, id int64, size int64) error {
	// update folder size
	query := `UPDATE folders SET size = size + $1, updated_at = NOW() WHERE id = $2`
	if _, err := tx.Exec(query, size, id); err != nil {
		return fmt.Errorf("failed to increase folder size: %w", err)
	}

	// propagate the size difference to parent folders
	err := r.updateParentFolderSizes(tx, id, size)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseFolderSize decreases a folder size and propagates the change to all parent folders
func (r *folderRepository) DecreaseFolderSize(tx *sql.Tx, id, size int64) error {
	query := `UPDATE folders SET size = size - $1, updated_at = NOW() WHERE id = $2`
	if _, err := tx.Exec(query, size, id); err != nil {
		return fmt.Errorf("failed to decrease folder size: %w", err)
	}

	if err := r.updateParentFolderSizes(tx, id, -size); err != nil {
		return err
	}
	return nil
}

// UpdateMultipleFoldersSize updates sizes of multiple folders using a batch update query
func (r *folderRepository) UpdateMultipleFoldersSize(tx *sql.Tx, folders []models.FolderSizeSimplified) error {
	if len(folders) == 0 {
		return nil
	}

	query := `
		UPDATE folders
		SET size = CASE id %s ELSE size END
		WHERE id IN (%s)
	`

	var sizeCases, idList string
	for i, folder := range folders {
		if i > 0 {
			idList += ", "
		}
		idList += fmt.Sprintf("%d", folder.ID)
		sizeCases += fmt.Sprintf(" WHEN %d THEN %d", folder.ID, folder.Size)
	}

	query = fmt.Sprintf(query, sizeCases, idList)
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to execute batch update: %w", err)
	}

	return nil
}

// propagates size changes to all parent folders recursively
func (r *folderRepository) updateParentFolderSizes(tx *sql.Tx, folderID, sizeDifference int64) error {
	for {
		var parentFolderID sql.NullInt64
		err := tx.QueryRow(`SELECT parent_folder_id FROM folders WHERE id = $1`, folderID).Scan(&parentFolderID)
		if err != nil {
			return fmt.Errorf("failed to retrieve parent folder ID: %w", err)
		}

		if !parentFolderID.Valid {
			break
		}

		query := `UPDATE folders SET size = size + $1, updated_at = NOW() WHERE id = $2`
		if _, err = tx.Exec(query, sizeDifference, parentFolderID.Int64); err != nil {
			return fmt.Errorf("failed to update parent folder size: %w", err)
		}

		folderID = parentFolderID.Int64
	}

	return nil
}

// recursiveDeleteFolder recursively deletes a folder and its subfolders
func (r *folderRepository) recursiveDeleteFolder(tx *sql.Tx, folderID int64) error {
	subfolders, err := r.GetFoldersInfo(folderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve subfolders for deletion: %w", err)
	}

	for _, subfolder := range subfolders {
		if subfolder.ID == folderID {
			continue
		}
		if err := r.recursiveDeleteFolder(tx, subfolder.ID); err != nil {
			return err
		}
	}

	query := `DELETE FROM folders WHERE id = $1`
	if _, err := tx.Exec(query, folderID); err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	return nil
}

// prevents creating a cycle in the folder hierarchy by ensuring no folder can be moved into one of its descendants
func (r *folderRepository) checkForCycle(tx *sql.Tx, folderID, newParentFolderID int64) error {
	currentID := newParentFolderID
	for currentID != 0 {
		if currentID == folderID {
			return errors.New("moving folder would create a cycle")
		}

		var parentID sql.NullInt64
		err := tx.QueryRow(`SELECT parent_folder_id FROM folders WHERE id = $1`, currentID).Scan(&parentID)
		if err != nil {
			return fmt.Errorf("failed to retrieve parent folder ID for cycle check: %w", err)
		}

		if !parentID.Valid {
			break
		}

		currentID = parentID.Int64
	}

	return nil
}
