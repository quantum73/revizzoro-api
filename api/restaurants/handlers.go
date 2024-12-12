package restaurants

import (
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
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

	rows, err := db.QueryContext(ctx, "SELECT * FROM restaurants")
	if err != nil {
		resp := network.NewInternalServerErrorResponse(network.InternalServerErrorBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}
	defer rows.Close()

	restaurants := make([]*model.Restaurant, 0)
	for rows.Next() {
		var (
			id         int
			name, link string
		)
		if err := rows.Scan(&id, &name, &link); err != nil {
			log.Errorf("Error scanning row: %s", err)
		}

		r, err := model.NewRestaurant(id, name, link)
		if err != nil {
			log.Errorf("Error creating restaurant sturcture: %s", err)
		}
		restaurants = append(restaurants, r)
	}

	if err := rows.Err(); err != nil {
		resp := network.NewInternalServerErrorResponse(network.InternalServerErrorBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	resp := network.NewSuccessDataResponse(network.OKBaseMessage, restaurants)
	ctx.JSON(resp.GetStatus(), resp)
}

func (c *controller) DetailByIdHandler(ctx *gin.Context) {
	db := c.db.GetInstance()

	restaurantId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		resp := network.NewBadRequestResponse(network.NotFoundBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	var (
		id         int
		name, link string
	)
	err = db.QueryRowContext(
		ctx, "SELECT * FROM restaurants AS r WHERE r.id = $1", restaurantId,
	).Scan(&id, &name, &link)
	if err != nil {
		log.Errorf("Error getting restaurant by `%d` id: %s", restaurantId, err)
		resp := network.NewBadRequestResponse(network.NotFoundBaseMessage)
		ctx.JSON(resp.GetStatus(), resp)
		return
	}

	restaurant, err := model.NewRestaurant(id, name, link)
	if err != nil {
		log.Errorf("Error creating restaurant sturcture: %s", err)
	}
	resp := network.NewSuccessDataResponse(network.OKBaseMessage, restaurant)
	ctx.JSON(resp.GetStatus(), resp)
}
