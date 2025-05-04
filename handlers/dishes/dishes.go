package dishes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quantum73/revizzoro-api/internal/common"
	"github.com/quantum73/revizzoro-api/internal/enums"
	"github.com/quantum73/revizzoro-api/internal/network"
	"github.com/quantum73/revizzoro-api/services/dishes"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var errDishIdMustBeInteger = errors.New("id must be integer")

type dishService interface {
	DetailById(ctx context.Context, idx int) (*dishes.DishReadDTO, error)
	List(ctx context.Context) ([]*dishes.DishReadDTO, error)
	Create(ctx context.Context, db *dishes.DishCreateDTO) (*dishes.DishReadDTO, error)
}

type DishController struct {
	database *sql.DB
	service  dishService
}

func NewDishController(database *sql.DB, service dishService) *DishController {
	return &DishController{database: database, service: service}
}

func (c *DishController) DetailById(w http.ResponseWriter, r *http.Request) {
	const op = "[handlers/dishes detailByIdHandler]"

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Warnf("%s incorrect `id` parameter in path: %s\n", op, err)
		network.BadRequestMessageResponse(w, errDishIdMustBeInteger.Error())
		return
	}

	dishDTO, err := c.service.DetailById(r.Context(), idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			network.NotFoundMessageResponse(w, fmt.Sprintf("dish by `%d` id not found", idInt))
		} else {
			network.ServerUnexpectedErrorMessageResponse(w)
		}
		return
	}

	marshalledDishDTO, err := json.Marshal(dishDTO)
	if err != nil {
		log.Warnf("%s error during marshall DTO to json: %s\n", op, err)
		network.ServerUnexpectedErrorMessageResponse(w)
		return
	}

	network.WriteResponse(
		w,
		http.StatusOK,
		enums.ContentTypeJSON,
		marshalledDishDTO,
	)
}

func (c *DishController) List(w http.ResponseWriter, r *http.Request) {
	const op = "[dishes listHandler]"

	dishesDTO, err := c.service.List(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			network.NotFoundMessageResponse(w, "there are no dishes")
		} else {
			network.ServerUnexpectedErrorMessageResponse(w)
		}
		return
	}

	dishesDTOAsBytes := make([][]byte, 0)
	for _, dish := range dishesDTO {
		marshalledDish, err := json.Marshal(dish)
		if err != nil {
			break
		}
		dishesDTOAsBytes = append(dishesDTOAsBytes, marshalledDish)
	}

	jsonPayload, err := common.ConcatenateBytesSlice(dishesDTOAsBytes)
	if err != nil {
		log.Warnf("%s error during marshall DTO to json: %s\n", op, err)
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

func (c *DishController) Create(w http.ResponseWriter, r *http.Request) {
	const op = "[dishes createHandler]"

	dishDTO := &dishes.DishCreateDTO{}
	err := json.NewDecoder(r.Body).Decode(dishDTO)
	if err != nil {
		log.Errorf("%s error during decoding request body to json: %s\n", op, err)
		network.BadRequestMessageResponse(w, err.Error())
		return
	}
	err = dishDTO.Validate()
	if err != nil {
		log.Errorf("%s dish DTO validation error: %s\n", op, err)
		network.BadRequestMessageResponse(w, err.Error())
		return
	}

	newDishDTO, err := c.service.Create(r.Context(), dishDTO)
	if err != nil {
		log.Warnf("%s error during creating dish: %s\n", op, err)
		network.ServerUnexpectedErrorMessageResponse(w)
		return
	}

	marshalledDishDTO, err := json.Marshal(newDishDTO)
	if err != nil {
		log.Warnf("%s error during marshall DTO to json: %s\n", op, err)
		network.ServerUnexpectedErrorMessageResponse(w)
		return
	}

	network.WriteResponse(
		w,
		http.StatusOK,
		enums.ContentTypeJSON,
		marshalledDishDTO,
	)
}
