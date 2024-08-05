package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// CompleteTransaction completes an ongoing transaction
// @Summary      Complete a transaction
// @Description  Completes the specified transaction by updating its status to "completed" and updating folder sizes
// @Tags         transaction
// @Param        user_id         header    int     true  "User ID"
// @Param        folder_id       path      int64   true  "Folder ID"
// @Param        transaction_id  path      int64   true  "Transaction ID"
// @Produce      json
// @Success      200  {object}  nil               "Transaction successfully completed"
// @Failure      400  {object}  ErrorResponse     "Invalid transaction_id"
// @Failure      404  {object}  ErrorResponse     "Transaction not found"
// @Failure      500  {object}  ErrorResponse     "Internal Server Error"
// @Router       /v1/folders/{folder_id}/transaction/{transaction_id}/complete [put]
func (h *Handler) CompleteTransaction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.completeTransaction(w, r)
	})
}

func (h *Handler) completeTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID, err := strconv.ParseInt(r.PathValue("transaction_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid transaction_id")
		return
	}

	// Get transaction from the db
	transaction, err := h.transactionService.GetTransactionByID(transactionID)
	if err != nil {
		FailedResponse(w, http.StatusInternalServerError, "Failed to get transaction")
		return
	}
	if transaction == nil {
		FailedResponse(w, http.StatusNotFound, fmt.Sprintf("Transaction '%d' not found", transactionID))
		return
	}

	// Get all parent folders affected by transaction
	allAffectedFolders, err := h.folderService.GetAllParentFolders(transaction.FolderID)
	if err != nil {
		log.Info().Msgf("Failed to get all parents for folder(%d): %s", transaction.FolderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Generate a list of keys to get all folder sizes
	keys := make([]string, len(allAffectedFolders))
	for i, id := range allAffectedFolders {
		keys[i] = fmt.Sprintf("folder:%d", id)
	}
	ctx := context.Background()
	// Get all folder sizes from the cache
	foldersData, err := h.rc.GetMultipleFolders(ctx, keys)
	if err != nil {
		log.Info().Msgf("Failed to get folders from cache(%d): %s", transaction.FolderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Update all folders sizes
	err = h.folderService.UpdateMultipleFoldersSize(foldersData)
	if err != nil {
		log.Info().Msgf("Failed to update all sizes(%d): %s", transaction.FolderID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Update all sizes in the database
	err = h.transactionService.UpdateTransactionStatus(transactionID, "completed")
	if err != nil {
		log.Info().Msgf("Failed to complete transaction(%d): %s", transactionID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to complete transaction")
		return
	}

	SuccessfulResponse(w, http.StatusOK, nil)
}
