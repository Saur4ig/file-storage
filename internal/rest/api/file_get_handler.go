package api

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// GetFile returns file data by it`s id
// @Summary      Get a file
// @Tags         file
// @Param        folder_id   path      int64  true  "Folder ID"
// @Param        file_id      path      int64  true  "File ID"
// @Param        user_id     header    int    true "User ID"
// @Produce      json
// @Success      204  {object}  nil   "No Content"
// @Failure      400  {object}  ErrorResponse "Invalid file_id"
// @Failure      500  {object}  ErrorResponse "Internal Server Error"
// @Router       /v1/folders/{folder_id}/files/{file_id} [get]
func (h *Handler) GetFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.getFile(w, r)
	})
}

func (h *Handler) getFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := strconv.ParseInt(r.PathValue("file_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid file_id")
		return
	}

	file, err := h.fileService.GetFile(fileID)
	if err != nil {
		log.Warn().Msgf("failed to get file(%d): %s", fileID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to get file")
		return
	}

	SuccessfulResponse(w, http.StatusNoContent, file)
}
