package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/saur4ig/file-storage/internal/models"
)

// GetFolder retrieves folder size information
// @Summary      Get folder information
// @Description  Retrieves size details of the folder specified by folder ID, similar to the 'du' command, but in JSON format
// @Tags         folder
// @Param        user_id   header    int     true  "User ID"
// @Param        folder_id path      int64   true  "Folder ID"
// @Produce      json
// @Success      200  {array}   Size             "Folder size information retrieved successfully"
// @Failure      400  {object}  ErrorResponse    "Invalid folder_id"
// @Failure      500  {object}  ErrorResponse    "Failed to get folder information"
// @Router       /v1/folders/{folder_id} [get]
func (h *Handler) GetFolder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.getFolder(w, r)
	})
}

// Size represents the structure for folder size details
type Size struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

func (h *Handler) getFolder(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(r.PathValue("folder_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid folder_id")
		return
	}

	data, err := h.folderService.GetFolderInfo(folderID)
	if err != nil {
		log.Info().Msgf("Failed to get folder info: %s", err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to get folder info")
		return
	}

	SuccessfulResponse(w, http.StatusOK, ListFolderSizes(data))
}

// ListFolderSizes converts each FolderSize to a human-readable format and returns a list of sizes
func ListFolderSizes(folders []models.FolderSize) []Size {
	var formattedSizes []Size
	for _, folder := range folders {
		formattedSizes = append(formattedSizes, Size{
			Name: folder.Name,
			Size: humanReadableSize(folder.Size),
		})
	}
	return formattedSizes
}

// humanReadableSize converts bytes to a human-readable format (KB, MB, GB, etc.)
func humanReadableSize(size int64) string {
	const (
		_          = iota // ignore first value by assigning to blank identifier
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
		EB
	)

	sizeFloat := float64(size)
	switch {
	case sizeFloat >= EB:
		return fmt.Sprintf("%.2f EB", sizeFloat/EB)
	case sizeFloat >= PB:
		return fmt.Sprintf("%.2f PB", sizeFloat/PB)
	case sizeFloat >= TB:
		return fmt.Sprintf("%.2f TB", sizeFloat/TB)
	case sizeFloat >= GB:
		return fmt.Sprintf("%.2f GB", sizeFloat/GB)
	case sizeFloat >= MB:
		return fmt.Sprintf("%.2f MB", sizeFloat/MB)
	case sizeFloat >= KB:
		return fmt.Sprintf("%.2f KB", sizeFloat/KB)
	default:
		return fmt.Sprintf("%d Bytes", size)
	}
}
