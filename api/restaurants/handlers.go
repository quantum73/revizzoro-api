package restaurants

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"strconv"
)

func ListHandler(ctx *gin.Context) {
	resp := network.NewSuccessDataResponse(network.OKMessage, gin.H{"restaurants": []string{}})
	ctx.JSON(resp.GetStatus(), resp)
}

func DetailByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		resp := network.NewBadRequestResponse("restaurant not found")
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	newRestaurant, err := model.NewRestaurant(idAsInt, "MockName", "http://some-rest.com")
	if err != nil {
		resp := network.NewBadRequestResponse("something wrong with restaurant object")
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKMessage, newRestaurant)
	ctx.JSON(resp.GetStatus(), resp)
}
