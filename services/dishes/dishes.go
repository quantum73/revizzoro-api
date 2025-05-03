package dishes

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type DishService struct {
	db *sql.DB
}

func NewDishService(database *sql.DB) *DishService {
	return &DishService{db: database}
}

func (d DishService) DetailById(ctx context.Context, idx int) (*DishDTO, error) {
	const op = "[services/dishes DetailById]"

	dish := dishDB{}
	row := d.db.QueryRowContext(ctx, "SELECT * FROM dishes WHERE id = $1", idx)
	err := row.Scan(&dish.ID, &dish.Name, &dish.Price, &dish.Score, &dish.RestaurantID)
	if err != nil {
		log.Warnf("%s Error getting dish by `%d` id from db: %v\n", op, idx, err)
		return &DishDTO{}, err
	}

	return dish.toDTO(), nil
}

func (d DishService) List(ctx context.Context) ([]*DishDTO, error) {
	const op = "[services/dishes List]"

	dishesDTO := make([]*DishDTO, 0)

	rows, err := d.db.QueryContext(ctx, "SELECT * FROM dishes")
	if err != nil {
		log.Warnf("%s Error query dishes: %v\n", op, err)
		return dishesDTO, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Warnf("%s Error closing rows: %v\n", op, err)
		}
	}()

	for rows.Next() {
		var dish dishDB
		err := rows.Scan(&dish.ID, &dish.Name, &dish.Price, &dish.Score, &dish.RestaurantID)
		if err != nil {
			log.Warnf("%s Error getting dishes: %v\n", op, err)
			return dishesDTO, err
		}

		dishesDTO = append(dishesDTO, dish.toDTO())
	}

	return dishesDTO, nil
}
