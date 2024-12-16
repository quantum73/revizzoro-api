package dishes

import (
	"context"
	"errors"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const notFoundMessage = "dish not found"
const badRequestMessage = "error during getting dishes"

type Service interface {
	GetAll() ([]model.Dish, network.ApiError)
	GetOneByID(id int) (*model.Dish, network.ApiError)
}

type service struct {
	db      postgres.Database
	context context.Context
}

func NewService(ctx context.Context, db postgres.Database) Service {
	return &service{context: ctx, db: db}
}

func (s *service) GetAll() ([]model.Dish, network.ApiError) {
	db := s.db.GetInstance()

	var dishes []model.Dish
	result := db.Find(&dishes)
	if err := result.Error; err != nil {
		log.Errorf("error during getting dishes: %s", result.Error)
		return nil, network.NewBadRequestError(badRequestMessage, err)
	}

	return dishes, nil
}

func (s *service) GetOneByID(dishId int) (*model.Dish, network.ApiError) {
	db := s.db.GetInstance()

	var dish model.Dish
	result := db.First(&dish, dishId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, network.NewNotFoundError(notFoundMessage, result.Error)
		} else {
			log.Errorf(
				"unexpected error during getting dish by `%d` id: %s",
				dishId, result.Error,
			)
			return nil, network.NewBadRequestError(badRequestMessage, result.Error)
		}
	}

	return &dish, nil
}
