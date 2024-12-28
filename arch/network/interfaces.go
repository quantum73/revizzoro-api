package network

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError interface {
	GetCode() int
	GetMessage() string
	Error() string
	Unwrap() error
}

type Response interface {
	GetStatus() int
	GetMessage() string
	GetData() any
}

type BaseController interface {
	MountRoutes(group *gin.RouterGroup)
}

type Dto[T any] interface {
	GetValue() *T
	ValidateErrors(errs validator.ValidationErrors) ([]string, error)
}
