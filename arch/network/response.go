package network

import (
	"net/http"
)

const (
	OKBaseMessage                  string = "OK"
	NotFoundBaseMessage            string = "Not found"
	InternalServerErrorBaseMessage string = "Something went wrong on server"
)

type response struct {
	Status  int    `json:"status" binding:"required"`
	Message string `json:"message" binding:"required"`
	Data    any    `json:"data,omitempty" binding:"required,omitempty"`
}

func (r *response) GetStatus() int {
	return r.Status
}

func (r *response) GetMessage() string {
	return r.Message
}

func (r *response) GetData() any {
	return r.Data
}

func NewSuccessDataResponse(message string, data any) Response {
	return &response{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

func NewSuccessMsgResponse(message string) Response {
	return &response{
		Status:  http.StatusOK,
		Message: message,
	}
}

func NewBadRequestResponse(message string) Response {
	return &response{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundResponse(message string) Response {
	return &response{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewInternalServerErrorResponse(message string) Response {
	return &response{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}
