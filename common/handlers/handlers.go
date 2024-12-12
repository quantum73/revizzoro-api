package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/arch/network"
)

func DefaultNotFoundHandler(ctx *gin.Context) {
	resp := network.NewNotFoundResponse(network.NotFoundBaseMessage)
	ctx.JSON(resp.GetStatus(), resp)
}
