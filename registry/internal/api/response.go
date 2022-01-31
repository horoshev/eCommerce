package api

import (
	"encoding/json"
	"net/http"
)

type ServerResponse struct {
	Name    string `json:"name" extensions:"x-order=0"`
	Version string `json:"version" extensions:"x-order=1"`
	Swagger string `json:"swagger" extensions:"x-order=2"`
}

type Response struct {
	Message string `json:"message" extensions:"x-order=0"`
}

func OkResponse(w http.ResponseWriter, body interface{}) {
	response(w, http.StatusOK, body)
}

func ErrorResponse(w http.ResponseWriter, err error) {
	BadRequestResponse(w, err.Error())
}

func UnauthorizedResponse(w http.ResponseWriter, message string) {
	response(w, http.StatusUnauthorized, Response{message})
}

func BadRequestResponse(w http.ResponseWriter, message string) {
	response(w, http.StatusBadRequest, Response{message})
}

func InternalErrorResponse(w http.ResponseWriter, message string) {
	response(w, http.StatusInternalServerError, Response{message})
}

func response(w http.ResponseWriter, code int, body interface{}) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
