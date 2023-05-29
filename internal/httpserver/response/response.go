package response

import (
	"encoding/json"
	"net/http"

	"github.com/defryheryanto/mini-wallet/internal/errors"
)

func Success(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   data,
	})
}

func Failed(w http.ResponseWriter, err error) {
	isHandledError, handledError := handleError(err)
	if isHandledError {
		failedHandled(w, handledError.HttpStatus, handledError.Data)
		return
	}

	failedError(w, err)
}

func failedError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": err.Error(),
	})
}

func failedHandled(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "fail",
		"data": map[string]interface{}{
			"error": data,
		},
	})
}

func handleError(err error) (bool, errors.HandledError) {
	isHandledError := false
	var definedError errors.HandledError

	switch err := err.(type) {
	case errors.HandledError:
		definedError = err
		isHandledError = true
	}

	return isHandledError, definedError
}
