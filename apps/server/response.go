package main

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Error responseError `json:"error"`
}

type responseError struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorBody{
		Error: responseError{Message: message},
	})
}

func writeResult(w http.ResponseWriter, value any, err error) {
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, value)
}
