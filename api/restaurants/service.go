package restaurants

import (
	"context"
	"errors"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const notFoundMessage = "restaurant not found"
const badRequestMessage = "error during getting restaurants"

type Service interface {
	GetAll() ([]model.Restaurant, network.ApiError)
	GetOneByID(id int) (*model.Restaurant, network.ApiError)
}

type service struct {
	db      postgres.Database
	context context.Context
}

func NewService(ctx context.Context, db postgres.Database) Service {
	return &service{context: ctx, db: db}
}

func (s *service) GetAll() ([]model.Restaurant, network.ApiError) {
	db := s.db.GetInstance()

	var restaurants []model.Restaurant
	result := db.Find(&restaurants)
	if err := result.Error; err != nil {
		log.Errorf("error during getting restaurants: %s", result.Error)
		return nil, network.NewBadRequestError(badRequestMessage, err)
	}

	return restaurants, nil
}

func (s *service) GetOneByID(restaurantId int) (*model.Restaurant, network.ApiError) {
	db := s.db.GetInstance()

	var restaurant model.Restaurant
	result := db.First(&restaurant, restaurantId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, network.NewNotFoundError(notFoundMessage, result.Error)
		} else {
			log.Errorf(
				"unexpected error during getting restaurant by `%d` id: %s",
				restaurantId, result.Error,
			)
			return nil, network.NewBadRequestError(badRequestMessage, result.Error)
		}
	}

	return &restaurant, nil
}
