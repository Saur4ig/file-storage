package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/saur4ig/file-storage/internal/rest/middleware"
)

// CreateFolder creates a new folder
// @Summary      Create a new folder
// @Description  Creates a new folder with the given name and parent folder ID for the authenticated user.
// @Tags         folder
// @Param        user_id         header    int         true  "User ID"
// @Param        requestData     body      RequestData true  "Folder creation request payload"
// @Produce      json
// @Success      201  {object}  NewFolderResponse "Folder successfully created"
// @Failure      400  {object}  ErrorResponse     "Invalid request data or folder creation failed"
// @Failure      500  {object}  ErrorResponse     "Internal Server Error"
// @Router       /v1/folders [post]
func (h *Handler) CreateFolder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.createFolder(w, r)
	})
}

// RequestData structure of the new folder request
type RequestData struct {
	Name           string `json:"name"`
	ParentFolderID int64  `json:"parent_folder_id"`
}

// NewFolderResponse response structure with new folder id
type NewFolderResponse struct {
	FolderID int64 `json:"folder_id"`
}

func (h *Handler) createFolder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDHeaderKey).(int)

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// Decode the JSON data into struct
	var data RequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Failed to decode request")
		return
	}

	// Create folder in database
	folderID, err := h.folderService.CreateFolder(userID, data.Name, data.ParentFolderID)
	if err != nil {
		log.Info().Msgf("Error on folder creation: %s", err.Error())
		FailedResponse(w, http.StatusBadRequest, "Failed to create folder")
		return
	}

	SuccessfulResponse(w, http.StatusCreated, NewFolderResponse{
		FolderID: folderID,
	})
}
