package api

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// DeleteFile removes user file by id
// @Summary      Remove a file
// @Description  Deletes a file from the specified folder by file ID.
// @Tags         file
// @Param        folder_id   path      int64  true  "Folder ID"
// @Param        file_id      path      int64  true  "File ID"
// @Param        user_id     header    int    true "User ID"
// @Produce      json
// @Success      204  {object}  nil   "No Content"
// @Failure      400  {object}  ErrorResponse "Invalid folder_id or file_id"
// @Failure      500  {object}  ErrorResponse "Internal Server Error"
// @Router       /v1/folders/{folder_id}/files/{file_id} [delete]
func (h *Handler) DeleteFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.deleteFile(w, r)
	})
}

func (h *Handler) deleteFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(r.PathValue("file_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid file_id")
		return
	}

	err = h.fileService.DeleteFile(fileID)
	if err != nil {
		log.Warn().Msgf("failed to remove file(%d): %s", fileID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	SuccessfulResponse(w, http.StatusNoContent, nil)
}
