package api

import (
	"encoding/json"
	"net/http"

	"github.com/ziollek/etherscription/pkg/logging"
)

func ErrorResponse(errorCode int, errorMsg string, w http.ResponseWriter) {
	w.WriteHeader(errorCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&JSONErrorResponse{Error: &GenericError{Status: errorCode, Title: errorMsg}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Response(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if response != nil {
		if err := encoder.Encode(response); err != nil {
			ErrorResponse(http.StatusInternalServerError, "Internal Server Error", w)
			logging.Logger().Err(err)
		}
	}
}
