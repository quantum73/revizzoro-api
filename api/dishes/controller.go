package dishes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

const (
	notFoundMessage   = "dish not found"
	badRequestMessage = "error during getting dishes"
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

	var dishes []model.Dish
	result := db.Find(&dishes)
	if err := result.Error; err != nil {
		log.Errorf("error during getting dishes: %s", result.Error)
		resp := network.NewBadRequestResponse(notFoundMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, dishes)
	ctx.JSON(resp.GetStatus(), resp)
}

func (c *controller) DetailByIdHandler(ctx *gin.Context) {
	dishId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := network.NewNotFoundResponse(notFoundMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	db := c.db.GetInstance()

	var dish model.Dish
	result := db.First(&dish, dishId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			resp := network.NewNotFoundResponse(notFoundMessage)
			ctx.JSON(resp.GetStatus(), resp)
			return
		} else {
			log.Errorf(
				"unexpected error during getting dish by `%d` id: %s",
				dishId, result.Error,
			)
			resp := network.NewNotFoundResponse(badRequestMessage)
			ctx.JSON(resp.GetStatus(), resp)
			return
		}
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, dish)
	ctx.JSON(resp.GetStatus(), resp)
}
