package api

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// RemoveFolder removes a folder and updates all related sizes
// @Summary      Remove a folder
// @Description  Deletes a specified folder and updates size calculations for all related folders
// @Tags         folder
// @Param        user_id     header    int     true  "User ID"
// @Param        folder_id   path      int64   true  "Folder ID"
// @Produce      json
// @Success      204  {object}  nil               "Folder successfully removed"
// @Failure      400  {object}  ErrorResponse     "Invalid folder_id"
// @Failure      500  {object}  ErrorResponse     "Failed to remove folder"
// @Router       /v1/folders/{folder_id} [delete]
func (h *Handler) RemoveFolder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.removeFolder(w, r)
	})
}

func (h *Handler) removeFolder(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(r.PathValue("folder_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid folder_id")
		return
	}

	// Remove folder in the database
	err = h.folderService.DeleteFolder(folderID)
	if err != nil {
		log.Info().Msgf("Failed to remove folder(%d): %s", folderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to remove folder")
		return
	}

	// Update size of all folders in cache

	SuccessfulResponse(w, http.StatusNoContent, nil)
}
