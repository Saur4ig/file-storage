package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// MoveFolder changes the parent folder of the requested folder and updates all sizes
// @Summary      Move a folder
// @Description  Changes the parent folder of a specified folder and updates the size calculations
// @Tags         folder
// @Param        user_id     header    int               true  "User ID"
// @Param        folder_id   path      int64             true  "Folder ID"
// @Param        moveFolder  body      MoveFolderRequest true  "New parent folder ID"
// @Produce      json
// @Success      200  {object}  nil               "Folder successfully moved"
// @Failure      400  {object}  ErrorResponse     "Invalid folder_id or request body"
// @Failure      500  {object}  ErrorResponse     "Internal Server Error"
// @Router       /v1/folders/{folder_id}/move [put]
func (h *Handler) MoveFolder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.moveFolder(w, r)
	})
}

// MoveFolderRequest structure to handle folder move requests
type MoveFolderRequest struct {
	NewFolderID int64 `json:"new_folder_id"`
}

func (h *Handler) moveFolder(w http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.ParseInt(r.PathValue("folder_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid folder_id")
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// Decode the JSON data into struct
	var data MoveFolderRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Failed to decode request")
		return
	}

	// Move folder and re-calculate sizes
	err = h.folderService.MoveFolder(folderID, data.NewFolderID)
	if err != nil {
		log.Info().Msgf("Failed to move folder(%d) to %d: %s", folderID, data.NewFolderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to move folder")
		return
	}

	SuccessfulResponse(w, http.StatusOK, nil)
}
