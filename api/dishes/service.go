package dishes

import (
	"context"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
	"github.com/quantum73/revizzoro-api/arch/network"
	"github.com/quantum73/revizzoro-api/arch/postgres"
	log "github.com/sirupsen/logrus"
)

const notFoundMessage = "dish not found"

type Service interface {
	GetAll() ([]*model.Dish, network.ApiError)
	GetOneByID(id int) (*model.Dish, network.ApiError)
}

type service struct {
	db      postgres.Database
	context context.Context
}

func NewService(ctx context.Context, db postgres.Database) Service {
	return &service{context: ctx, db: db.GetInstance()}
}

func (s *service) GetAll() ([]*model.Dish, network.ApiError) {
	db := s.db.GetInstance()

	dishes := make([]*model.Dish, 0)
	rows, err := db.QueryContext(s.context, "SELECT * FROM dishes")
	if err != nil {
		return dishes, network.NewNotFoundError(notFoundMessage, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, price, score, restaurantId int
			name                           string
		)
		if err := rows.Scan(&id, &name, &price, &score, &restaurantId); err != nil {
			log.Errorf("Error scanning row: %s", err)
			continue
		}

		r, err := model.NewDish(id, name, price, score, restaurantId)
		if err != nil {
			log.Errorf("Error creating dish sturcture: %s", err)
			continue
		}
		dishes = append(dishes, r)
	}

	if err := rows.Err(); err != nil {
		return dishes, network.NewNotFoundError(notFoundMessage, err)
	}

	return dishes, nil
}

func (s *service) GetOneByID(dishId int) (*model.Dish, network.ApiError) {
	db := s.db.GetInstance()

	var (
		id, price, score, restaurantId int
		name                           string
	)
	err := db.QueryRowContext(
		s.context,
		"SELECT * FROM dishes AS d WHERE d.id = $1",
		dishId,
	).Scan(&id, &name, &price, &score, &restaurantId)
	if err != nil {
		log.Errorf("Error getting dish by `%d` id: %s", dishId, err)
		return nil, network.NewNotFoundError(notFoundMessage, err)
	}

	dish, err := model.NewDish(id, name, price, score, restaurantId)
	if err != nil {
		log.Errorf("Error creating dish structure: %s", err)
		return nil, network.NewNotFoundError(notFoundMessage, err)
	}

	return dish, nil
}
