package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/saur4ig/file-storage/internal/rest/middleware"
)

// UploadFile uploads a file to S3 and saves the metadata in the database
// @Summary      Upload a file
// @Description  Uploads a file to S3 storage and saves the file details in the database. It also updates the folder size cache in there is no transaction
// @Tags         file
// @Param        user_id         header    int     true  "User ID"
// @Param        folder_id       path      int64   true  "Folder ID"
// @Param        file             formData  file     true  "File to upload"
// @Accept       multipart/form-data
// @Produce      json
// @Success      201  {object}  nil                 "File successfully uploaded"
// @Failure      400  {object}  ErrorResponse       "Invalid input parameters or file upload failed"
// @Failure      500  {object}  ErrorResponse       "Internal Server Error"
// @Router       /v1/folders/{folder_id}/files [post]
func (h *Handler) UploadFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.uploadFile(w, r)
	})
}

func (h *Handler) uploadFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDHeaderKey).(int)

	folderID, err := strconv.ParseInt(r.PathValue("folder_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid folder_id")
		return
	}

	// Get transaction from headers
	var transactionID *int64
	transactionIDStr := r.Header.Get("transaction_id")
	if transactionIDStr != "" {
		id, err := strconv.ParseInt(transactionIDStr, 10, 64)
		if err == nil {
			transactionID = &id
		}
	}

	// Get file
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Info().Msgf("Failed to get file from request: %s", err.Error())
		FailedResponse(w, http.StatusBadRequest, "Error occurred on file processing")
		return
	}
	defer file.Close()

	name := header.Filename
	size := header.Size // in bytes

	// Always nil - because s3 real logic not implemented
	fileURL, err := h.s3.UploadFile(nil, name)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Error occurred on file saving")
		return
	}

	// Save file in db and update
	err = h.fileService.UploadFile(folderID, userID, name, fileURL, size, transactionID)
	if err != nil {
		log.Info().Msgf("Failed to save file to db: %s", err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Error occurred on file saving")
		return
	}

	// Update folder cache
	ctx := context.Background()
	err = h.rc.SetOrUpdateFolderSize(ctx, folderID, size)
	if err != nil {
		log.Info().Msgf("Failed to save file to cache: %s", err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Error occurred on file caching")
		return
	}

	SuccessfulResponse(w, http.StatusCreated, nil)
}
