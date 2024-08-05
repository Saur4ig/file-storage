package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// MoveFile moves a file to a new folder
// @Summary      Move a file
// @Description  Moves a specified file to a new folder based on provided folder ID.
// @Tags         file
// @Param        folder_id   path      int64  true  "Folder ID"
// @Param        file_id      path      int64  true  "File ID"
// @Param        user_id     header    int    true  "User ID"
// @Param        moveFile    body      MoveFileRequest true  "Request payload containing the new folder ID"
// @Produce      json
// @Success      200  {object}  nil   "File successfully moved"
// @Failure      400  {object}  ErrorResponse "Invalid input parameters"
// @Failure      500  {object}  ErrorResponse "Internal Server Error"
// @Router       /v1/folders/{folder_id}/files/{file_id}/move [put]
func (h *Handler) MoveFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.moveFile(w, r)
	})
}

// MoveFileRequest represents the request payload to move a file
type MoveFileRequest struct {
	NewFolderID int64 `json:"new_folder_id"`
}

func (h *Handler) moveFile(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(r.PathValue("folder_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid folder_id")
		return
	}

	fileID, err := strconv.ParseInt(r.PathValue("file_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid file_id")
		return
	}

	// read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Info().Msgf("Failed to read body: %s", err.Error())
		FailedResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// decode the JSON data into struct
	var data MoveFileRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Failed to decode request")
		return
	}

	// move file and re-calculate sizes
	err = h.fileService.MoveFile(fileID, folderID, data.NewFolderID)
	if err != nil {
		log.Warn().Msgf("Failed to move file(%d) to %d: %s", fileID, data.NewFolderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to move file")
		return
	}

	SuccessfulResponse(w, http.StatusOK, nil)
}
