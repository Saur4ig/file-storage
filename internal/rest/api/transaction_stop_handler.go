package api

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// StopTransaction stops an ongoing transaction
// @Summary      Stop a transaction
// @Description  Stops the specified transaction by updating its status to "failed"
// @Tags         transaction
// @Param        user_id         header    int     true  "User ID"
// @Param        folder_id       path      int64   true  "Folder ID"
// @Param        transaction_id  path      int64   true  "Transaction ID"
// @Produce      json
// @Success      200  {object}  nil               "Transaction successfully stopped"
// @Failure      400  {object}  ErrorResponse     "Invalid transaction_id"
// @Failure      500  {object}  ErrorResponse     "Failed to stop transaction"
// @Router       /v1/folders/{folder_id}/transaction/{transaction_id}/stop [put]
func (h *Handler) StopTransaction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.stopTransaction(w, r)
	})
}

func (h *Handler) stopTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID, err := strconv.ParseInt(r.PathValue("transaction_id"), 10, 64)
	if err != nil {
		FailedResponse(w, http.StatusBadRequest, "Invalid transaction_id")
		return
	}

	err = h.transactionService.UpdateTransactionStatus(transactionID, "failed")
	if err != nil {
		log.Info().Msgf("Failed to stop transaction(%d): %s", transactionID, err.Error())
		FailedResponse(w, http.StatusInternalServerError, "Failed to stop transaction")
		return
	}

	SuccessfulResponse(w, http.StatusOK, nil)
}
