package dishes

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	log "github.com/sirupsen/logrus"
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
	db := c.db.GetInstance()

	rows, err := db.QueryContext(ctx, "SELECT * FROM dishes")
	if err != nil {
		resp := network.NewInternalServerErrorResponse(network.InternalServerErrorBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}
	defer rows.Close()

	dishes := make([]*model.Dish, 0)
	for rows.Next() {
		var (
			id, price, score, restaurantId int
			name                           string
		)
		if err := rows.Scan(&id, &name, &price, &score, &restaurantId); err != nil {
			log.Errorf("Error scanning row: %s", err)
		}

		r, err := model.NewDish(id, name, price, score, restaurantId)
		if err != nil {
			log.Errorf("Error creating dish structure: %s", err)
			continue
		}
		dishes = append(dishes, r)
	}

	if err := rows.Err(); err != nil {
		resp := network.NewInternalServerErrorResponse(network.InternalServerErrorBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, dishes)
	ctx.JSON(resp.GetStatus(), resp)
}

func (c *controller) DetailByIdHandler(ctx *gin.Context) {
	db := c.db.GetInstance()

	dishId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := network.NewBadRequestResponse(network.NotFoundBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	var (
		id, price, score, restaurantId int
		name                           string
	)
	err = db.QueryRowContext(
		ctx,
		"SELECT * FROM dishes AS d WHERE d.id = $1",
		dishId,
	).Scan(&id, &name, &price, &score, &restaurantId)
	if err != nil {
		log.Errorf("Error getting dish by `%d` id: %s", dishId, err)
		resp := network.NewBadRequestResponse(network.NotFoundBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	dish, err := model.NewDish(id, name, price, score, restaurantId)
	if err != nil {
		log.Errorf("Error creating dish structure: %s", err)
	}
	resp := network.NewSuccessDataResponse(network.OKBaseMessage, dish)
	ctx.JSON(resp.GetStatus(), resp)
}
