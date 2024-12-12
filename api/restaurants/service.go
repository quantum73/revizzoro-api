package restaurants

import (
	"context"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	log "github.com/sirupsen/logrus"
)

const notFoundMessage = "restaurant not found"

type Service interface {
	GetAll() ([]*model.Restaurant, network.ApiError)
	GetOneByID(id int) (*model.Restaurant, network.ApiError)
}

type service struct {
	db      postgres.Database
	context context.Context
}

func NewService(ctx context.Context, db postgres.Database) Service {
	return &service{context: ctx, db: db.GetInstance()}
}

func (s *service) GetAll() ([]*model.Restaurant, network.ApiError) {
	db := s.db.GetInstance()

	restaurants := make([]*model.Restaurant, 0)
	rows, err := db.QueryContext(s.context, "SELECT * FROM restaurants")
	if err != nil {
		return restaurants, network.NewNotFoundError(notFoundMessage, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         int
			name, link string
		)
		if err := rows.Scan(&id, &name, &link); err != nil {
			log.Errorf("Error scanning row: %s", err)
			continue
		}

		r, err := model.NewRestaurant(id, name, link)
		if err != nil {
			log.Errorf("Error creating restaurant sturcture: %s", err)
			continue
		}
		restaurants = append(restaurants, r)
	}

	if err := rows.Err(); err != nil {
		return restaurants, network.NewNotFoundError(notFoundMessage, err)
	}

	return restaurants, nil
}

func (s *service) GetOneByID(restaurantId int) (*model.Restaurant, network.ApiError) {
	db := s.db.GetInstance()

	var (
		id         int
		name, link string
	)
	err := db.QueryRowContext(
		s.context,
		"SELECT * FROM restaurants AS r WHERE r.id = $1",
		restaurantId,
	).Scan(&id, &name, &link)
	if err != nil {
		log.Errorf("Error getting restaurant by `%d` id: %s", restaurantId, err)
		return nil, network.NewNotFoundError(notFoundMessage, err)
	}

	restaurant, err := model.NewRestaurant(id, name, link)
	if err != nil {
		log.Errorf("Error creating restaurant sturcture: %s", err)
		return nil, network.NewNotFoundError(notFoundMessage, err)
	}
	return restaurant, nil
}
