package dishes

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
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
	d, _ := model.NewDish(1, "Example", 1000, 5, 1)
	resp := network.NewSuccessDataResponse(network.OKMessage, []*model.Dish{d})
	ctx.JSON(resp.GetStatus(), resp)
}

func (c *controller) DetailByIdHandler(ctx *gin.Context) {
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
