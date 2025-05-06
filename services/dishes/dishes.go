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

func (d DishService) DetailById(ctx context.Context, idx int) (*DishReadDTO, error) {
	const op = "[services/dishes DetailById]"

	var dish DishDB
	row := d.db.QueryRowContext(ctx, "SELECT * FROM dishes WHERE id = $1", idx)
	err := row.Scan(&dish.ID, &dish.Name, &dish.Price, &dish.Score, &dish.RestaurantID)
	if err != nil {
		log.Warnf("%s Error getting dish by `%d` id from db: %v\n", op, idx, err)
		return &DishReadDTO{}, err
	}

	dto, err := dish.ToReadDTO()
	if err != nil {
		log.Warnf("%s Error converting dish to DTO: %v\n", op, err)
		return &DishReadDTO{}, err
	}
	return dto, nil
}

func (d DishService) List(
	ctx context.Context, limit, offset uint,
) ([]*DishReadDTO, error) {
	const op = "[services/dishes List]"

	dishesDTO := make([]*DishReadDTO, 0)

	rows, err := d.db.QueryContext(
		ctx,
		"SELECT * FROM dishes LIMIT $1 OFFSET $2",
		limit, offset,
	)
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
		var dish DishDB
		err := rows.Scan(&dish.ID, &dish.Name, &dish.Price, &dish.Score, &dish.RestaurantID)
		if err != nil {
			log.Warnf("%s Error getting dishes: %v\n", op, err)
			return dishesDTO, err
		}

		dto, err := dish.ToReadDTO()
		if err != nil {
			log.Warnf("%s Error converting dish to DTO: %v\n", op, err)
			return dishesDTO, err
		}

		dishesDTO = append(dishesDTO, dto)
	}

	return dishesDTO, nil
}

func (d DishService) Create(ctx context.Context, dishDTO *DishCreateDTO) (*DishReadDTO, error) {
	const op = "[services/dishes Create]"

	dish, err := dishDTO.ToDB()
	if err != nil {
		log.Warnf("%s Error converting dish to DB: %v\n", op, err)
		return &DishReadDTO{}, err
	}

	var newDishId uint
	err = d.db.QueryRowContext(
		ctx,
		"INSERT INTO dishes (name, price, score, restaurant_id) VALUES ($1, $2, $3, $4) RETURNING id",
		dish.Name, dish.Price, dish.Score, dish.RestaurantID,
	).Scan(&newDishId)
	if err != nil {
		log.Errorf("%s Error inserting dish to db: %v\n", op, err)
		return &DishReadDTO{}, err
	}

	newDishDTO := &DishReadDTO{
		ID:           newDishId,
		Name:         dish.Name,
		Price:        dish.Price,
		Score:        dish.Score,
		RestaurantID: dish.RestaurantID,
	}
	return newDishDTO, nil
}
