package http_utils

import (
	"encoding/json"
	"github.com/guebu/common-utils/errors"
	"net/http"
)

//ToDo: Merge this utils file with utils in common_utils library...

func RespondBodyAsJson(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func RespondErrorAsJson(w http.ResponseWriter, err *errors.ApplicationError) {
	RespondBodyAsJson(w, err.AStatusCode, err)
}
