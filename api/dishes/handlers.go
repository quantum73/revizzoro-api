package dishes

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
	"net/http"
	"strconv"
)

func ListHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"dishes": []string{}})
}

func DetailByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "dish not found"})
		return
	}

	newDish, err := model.NewDish(idAsInt, "Long Bull", 1500, 5, 1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "something wrong with dish object"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"dish": newDish})
}
