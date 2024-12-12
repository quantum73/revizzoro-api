package restaurants

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	"strconv"
)

type controller struct {
	context context.Context
	service Service
}

func NewController(ctx context.Context, db postgres.Database) network.BaseController {
	return &controller{context: ctx, service: NewService(ctx, db)}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/:id", c.DetailByIdHandler)
	group.GET("/", c.ListHandler)
}

func (c *controller) ListHandler(ctx *gin.Context) {
	restaurants, apiErr := c.service.GetAll()
	if apiErr != nil {
		resp := network.NewBadRequestResponse(apiErr.GetMessage())
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, restaurants)
	ctx.JSON(resp.GetStatus(), resp)
}

func (c *controller) DetailByIdHandler(ctx *gin.Context) {
	restaurantId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := network.NewNotFoundResponse(network.NotFoundBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	restaurant, apiErr := c.service.GetOneByID(restaurantId)
	if apiErr != nil {
		resp := network.NewNotFoundResponse(apiErr.GetMessage())
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, restaurant)
	ctx.JSON(resp.GetStatus(), resp)
}
