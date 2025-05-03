package dishes

import (
	"context"
	"database/sql"
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
	DetailById(ctx context.Context, idx int) (*dishes.DishDTO, error)
	List(ctx context.Context) ([]*dishes.DishDTO, error)
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
		log.Warnf("%s incorrect `id` parameter in path: %s", op, err)
		network.MessageJSONResponse(w, http.StatusBadRequest, errDishIdMustBeInteger.Error())
		return
	}

	dishDTO, err := c.service.DetailById(r.Context(), idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			network.MessageJSONResponse(
				w,
				http.StatusNotFound,
				fmt.Sprintf("dish by `%d` id not found", idInt),
			)
		} else {
			network.MessageJSONResponse(
				w,
				http.StatusInternalServerError,
				"unexpected error",
			)
		}

		return
	}

	dishDTOPayload, err := dishDTO.MarshallJSON()
	if err != nil {
		log.Warnf("%s error during marshall DTO to json: %s", op, err)
		network.MessageJSONResponse(
			w,
			http.StatusInternalServerError,
			"unexpected error",
		)
		return
	}

	network.WriteResponse(
		w,
		http.StatusOK,
		enums.ContentTypeJSON,
		dishDTOPayload,
	)
}

func (c *DishController) List(w http.ResponseWriter, r *http.Request) {
	const op = "dishes listHandler"

	dishesDTO, err := c.service.List(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			network.MessageJSONResponse(
				w,
				http.StatusNotFound,
				"there are no dishes",
			)
		} else {
			network.MessageJSONResponse(
				w,
				http.StatusInternalServerError,
				"unexpected error",
			)
		}

		return
	}

	dishesDTOAsBytes := make([][]byte, 0)
	for _, dish := range dishesDTO {
		p, err := dish.MarshallJSON()
		if err != nil {
			break
		}
		dishesDTOAsBytes = append(dishesDTOAsBytes, p)
	}

	jsonPayload, err := common.ConcatenateBytesSlice(dishesDTOAsBytes)
	if err != nil {
		log.Warnf("%s error during marshall DTO to json: %s\n", op, err)
		network.MessageJSONResponse(
			w,
			http.StatusInternalServerError,
			"unexpected error",
		)
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
	const op = "dishes createHandler"
	log.Infof("%s executed", op)

	network.MessageJSONResponse(w, http.StatusOK, "not implemented")
}
