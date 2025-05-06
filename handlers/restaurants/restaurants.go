package restaurants

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quantum73/revizzoro-api/internal/common"
	"github.com/quantum73/revizzoro-api/internal/enums"
	"github.com/quantum73/revizzoro-api/internal/network"
	"github.com/quantum73/revizzoro-api/internal/pagination"
	"github.com/quantum73/revizzoro-api/services/restaurants"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var errRestaurantIdMustBeInteger = errors.New("id must be integer")

type restaurantService interface {
	DetailById(ctx context.Context, idx int) (*restaurants.RestaurantReadDTO, error)
	List(ctx context.Context, limit, offset uint) ([]*restaurants.RestaurantReadDTO, error)
	Create(ctx context.Context, restaurantDTO *restaurants.RestaurantCreateDTO) (*restaurants.RestaurantReadDTO, error)
}

type RestaurantController struct {
	database  *sql.DB
	service   restaurantService
	paginator *pagination.Paginator
}

func NewRestaurantController(
	database *sql.DB, service restaurantService, paginator *pagination.Paginator,
) *RestaurantController {
	return &RestaurantController{database: database, service: service, paginator: paginator}
}

func (c *RestaurantController) DetailById(w http.ResponseWriter, r *http.Request) {
	const op = "[handlers/restaurants detailByIdHandler]"

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Warnf("%s incorrect `id` parameter in path: %s\n", op, err)
		network.BadRequestMessageResponse(w, errRestaurantIdMustBeInteger.Error())
		return
	}

	restaurantDTO, err := c.service.DetailById(r.Context(), idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			network.NotFoundMessageResponse(
				w,
				fmt.Sprintf("restaurant by `%d` id not found", idInt),
			)
		} else {
			network.ServerUnexpectedErrorMessageResponse(w)
		}
		return
	}

	marshalledRestaurantDTO, err := json.Marshal(restaurantDTO)
	if err != nil {
		log.Warnf("%s error during marshall restaurant DTO to json: %s\n", op, err)
		network.ServerUnexpectedErrorMessageResponse(w)
		return
	}

	network.WriteResponse(
		w,
		http.StatusOK,
		enums.ContentTypeJSON,
		marshalledRestaurantDTO,
	)
}

func (c *RestaurantController) List(w http.ResponseWriter, r *http.Request) {
	const op = "[handlers/restaurants listHandler]"

	pageLimit := c.paginator.LimitFromQueryParams(r.URL.Query())
	pageOffset := c.paginator.OffsetFromQueryParams(r.URL.Query())
	log.Infof("%s limit=%d offset=%d\n", op, pageLimit, pageOffset)

	restaurantsDTO, err := c.service.List(r.Context(), pageLimit, pageOffset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			network.NotFoundMessageResponse(w, "there are no restaurants")
		} else {
			network.ServerUnexpectedErrorMessageResponse(w)
		}
		return
	}

	restaurantsDTOAsBytes := make([][]byte, 0)
	for _, restaurant := range restaurantsDTO {
		marshalledRestaurant, err := json.Marshal(restaurant)
		if err != nil {
			break
		}
		restaurantsDTOAsBytes = append(restaurantsDTOAsBytes, marshalledRestaurant)
	}

	jsonPayload, err := common.ConcatenateBytesSlice(restaurantsDTOAsBytes)
	if err != nil {
		log.Warnf("%s error during marshall restaurants DTO to json: %s\n", op, err)
		network.ServerUnexpectedErrorMessageResponse(w)
		return
	}

	network.WriteResponse(
		w,
		http.StatusOK,
		enums.ContentTypeJSON,
		jsonPayload,
	)
}

func (c *RestaurantController) Create(w http.ResponseWriter, r *http.Request) {
	const op = "[handlers/restaurants createHandler]"

	restaurantDTO := &restaurants.RestaurantCreateDTO{}
	err := json.NewDecoder(r.Body).Decode(restaurantDTO)
	if err != nil {
		log.Errorf("%s error during decoding request body to json: %s\n", op, err)
		network.BadRequestMessageResponse(w, err.Error())
		return
	}
	err = restaurantDTO.Validate()
	if err != nil {
		log.Errorf("%s restaurant DTO validation error: %s\n", op, err)
		network.BadRequestMessageResponse(w, err.Error())
		return
	}

	newRestaurantDTO, err := c.service.Create(r.Context(), restaurantDTO)
	if err != nil {
		log.Warnf("%s error during creating restaurant: %s\n", op, err)
		network.ServerUnexpectedErrorMessageResponse(w)
		return
	}

	marshalledRestaurantDTO, err := json.Marshal(newRestaurantDTO)
	if err != nil {
		log.Warnf("%s error during marshall restaurant DTO to json: %s\n", op, err)
		network.ServerUnexpectedErrorMessageResponse(w)
		return
	}

	network.WriteResponse(
		w,
		http.StatusOK,
		enums.ContentTypeJSON,
		marshalledRestaurantDTO,
	)
}
