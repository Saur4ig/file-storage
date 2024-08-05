package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	error string
}

func sendResponse(w http.ResponseWriter, httpStatus int, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		log.Warn().Msgf("Error occured during marshaling the response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, err = w.Write(data)
	if err != nil {
		log.Warn().Msgf("Error occured during the response: %s", err.Error())
	}
}

func SuccessfulResponse(w http.ResponseWriter, httpStatus int, body interface{}) {
	sendResponse(w, httpStatus, body)
}

func FailedResponse(w http.ResponseWriter, httpStatus int, message string) {
	sendResponse(w, httpStatus, ErrorResponse{
		error: "message",
	})
}
