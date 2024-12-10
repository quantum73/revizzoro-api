package restaurants

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
	"net/http"
	"strconv"
)

func ListHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"restaurants": []string{}})
}

func DetailByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "restaurant not found"})
		return
	}

	newRestaurant, err := model.NewRestaurant(idAsInt, "MockName", "http://some-rest.com")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "something wrong with restaurant object"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"restaurant": newRestaurant})
}
