package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/saur4ig/file-storage/internal/models"
	"github.com/saur4ig/file-storage/internal/rest/middleware"
)

// StartTransaction begins a new transaction for a specified folder
// @Summary      Start a new transaction
// @Description  Initiates a new transaction for the specified folder
// @Tags         transaction
// @Param        user_id     header    int     true  "User ID"
// @Param        folder_id   path      int64   true  "Folder ID"
// @Produce      json
// @Success      201  {object}  TransactionStartResponse "Transaction successfully started"
// @Failure      400  {object}  ErrorResponse            "Invalid folder_id"
// @Failure      500  {object}  ErrorResponse            "Internal Server Error"
// @Router       /v1/folders/{folder_id}/transaction/start [post]
func (h *Handler) StartTransaction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.startTransaction(w, r)
	})
}

// TransactionStartResponse structure to respond with transaction ID
type TransactionStartResponse struct {
	TransactionID int64 `json:"transaction_id"`
}

func (h *Handler) startTransaction(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDHeaderKey).(int)

	folderID, err := strconv.ParseInt(r.PathValue("folder_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid folder_id")
		return
	}

	transactionID, err := h.transactionService.CreateTransaction(userID, folderID)
	if err != nil {
		log.Info().Msgf("Failed to create transaction for folder(%d): %s", folderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Get all parent folders whose size will be affected by the transaction
	allAffectedFolders, err := h.folderService.GetAllParentFolders(folderID)
	if err != nil {
		log.Info().Msgf("Failed to get all parents for folder(%d): %s", folderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Ensure that all of them are in cache
	go func(folders []models.FolderSizeSimplified) {
		ctx := context.Background()
		for _, folder := range folders {
			err := h.rc.SetOrUpdateFolderSize(ctx, folder.ID, folder.Size)
			if err != nil {
				log.Warn().Msgf("Error on updating folder(%d) size: %s", folder.ID, err.Error())
			}
		}
	}(allAffectedFolders)

	SuccessfulResponse(w, http.StatusCreated, TransactionStartResponse{
		TransactionID: transactionID,
	})
}
