package restaurants

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

const (
	notFoundMessage   = "restaurant not found"
	badRequestMessage = "error during getting restaurants"
)

type controller struct {
	context context.Context
	db      postgres.Database
}

func NewController(ctx context.Context, db postgres.Database) network.BaseController {
	return &controller{context: ctx, db: db}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/:id", c.DetailByIdHandler)
	group.GET("/", c.ListHandler)
}

func (c *controller) ListHandler(ctx *gin.Context) {
	db := c.db.GetInstance()

	var restaurants []model.Restaurant
	result := db.Find(&restaurants)
	if err := result.Error; err != nil {
		log.Errorf("error during getting restaurants: %s", result.Error)
		resp := network.NewBadRequestResponse(badRequestMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, restaurants)
	ctx.JSON(resp.GetStatus(), resp)
}

func (c *controller) DetailByIdHandler(ctx *gin.Context) {
	restaurantId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := network.NewNotFoundResponse(notFoundMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	db := c.db.GetInstance()

	var restaurant model.Restaurant
	result := db.First(&restaurant, restaurantId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			resp := network.NewNotFoundResponse(notFoundMessage)
			ctx.JSON(resp.GetStatus(), resp)
			return
		} else {
			log.Errorf(
				"unexpected error during getting restaurant by `%d` id: %s",
				restaurantId, result.Error,
			)
			resp := network.NewNotFoundResponse(badRequestMessage)
			ctx.JSON(resp.GetStatus(), resp)
			return
		}
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, restaurant)
	ctx.JSON(resp.GetStatus(), resp)
}
