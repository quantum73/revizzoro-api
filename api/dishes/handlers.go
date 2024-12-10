package dishes

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"strconv"
)

func ListHandler(ctx *gin.Context) {
	resp := network.NewSuccessDataResponse(network.OKMessage, gin.H{"dishes": []string{}})
	ctx.JSON(resp.GetStatus(), resp)
}

func DetailByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		resp := network.NewBadRequestResponse("dish not found")
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	newDish, err := model.NewDish(idAsInt, "Long Bull", 1500, 5, 1)
	if err != nil {
		resp := network.NewBadRequestResponse("something wrong with dish object")
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKMessage, newDish)
	ctx.JSON(resp.GetStatus(), resp)
}
