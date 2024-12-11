package restaurants

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	"strconv"
)

type controller struct {
	db postgres.Database
}

func NewController(db postgres.Database) network.BaseController {
	return &controller{db: db}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/:id", c.DetailByIdHandler)
	group.GET("/", c.ListHandler)
}

func (c *controller) ListHandler(ctx *gin.Context) {
	r, _ := model.NewRestaurant(1, "Example", "https://example-rest.com")
	resp := network.NewSuccessDataResponse(network.OKMessage, []*model.Restaurant{r})
	ctx.JSON(resp.GetStatus(), resp)
}

func (c *controller) DetailByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		resp := network.NewBadRequestResponse("restaurant not found")
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	newRestaurant, err := model.NewRestaurant(idAsInt, "MockName", "https://some-rest.com")
	if err != nil {
		resp := network.NewBadRequestResponse("something wrong with restaurant object")
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKMessage, newRestaurant)
	ctx.JSON(resp.GetStatus(), resp)
}
