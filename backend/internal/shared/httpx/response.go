package httpx

import (
	"encoding/json"
	"net/http"
)

type SuccessEnvelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorEnvelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, SuccessEnvelope{Success: true, Data: data})
}

func Error(w http.ResponseWriter, status int, message string, errors interface{}) {
	JSON(w, status, ErrorEnvelope{Success: false, Message: message, Errors: errors})
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, SuccessEnvelope{Success: true, Data: data})
}

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
