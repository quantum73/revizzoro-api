package network

import "github.com/gin-gonic/gin"

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
