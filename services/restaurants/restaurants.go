package restaurants

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type RestaurantService struct {
	db *sql.DB
}

func NewRestaurantService(database *sql.DB) *RestaurantService {
	return &RestaurantService{db: database}
}

func (r RestaurantService) DetailById(ctx context.Context, idx int) (*RestaurantReadDTO, error) {
	const op = "[services/restaurants DetailById]"

	var restaurant RestaurantDB
	row := r.db.QueryRowContext(ctx, "SELECT * FROM restaurants WHERE id = $1", idx)
	err := row.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Link)
	if err != nil {
		log.Warnf("%s Error getting restaurant by `%d` id from db: %v\n", op, idx, err)
		return &RestaurantReadDTO{}, err
	}

	dto, err := restaurant.ToReadDTO()
	if err != nil {
		log.Warnf("%s Error converting restaurant to DTO: %v\n", op, err)
		return &RestaurantReadDTO{}, err
	}
	return dto, nil
}

func (r RestaurantService) List(
	ctx context.Context, limit, offset uint,
) ([]*RestaurantReadDTO, error) {
	const op = "[services/restaurants DetailById]"

	restaurantsDTO := make([]*RestaurantReadDTO, 0)

	rows, err := r.db.QueryContext(
		ctx,
		"SELECT * FROM restaurants LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		log.Warnf("%s Error query restaurants: %v\n", op, err)
		return restaurantsDTO, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Warnf("%s Error closing rows: %v\n", op, err)
		}
	}()

	for rows.Next() {
		var restaurant RestaurantDB
		err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Link)
		if err != nil {
			log.Warnf("%s Error getting restaurants: %v\n", op, err)
			return restaurantsDTO, err
		}

		dto, err := restaurant.ToReadDTO()
		if err != nil {
			log.Warnf("%s Error converting restaurant to DTO: %v\n", op, err)
			return restaurantsDTO, err
		}

		restaurantsDTO = append(restaurantsDTO, dto)
	}

	return restaurantsDTO, nil
}

func (r RestaurantService) Create(
	ctx context.Context, restaurantDTO *RestaurantCreateDTO,
) (*RestaurantReadDTO, error) {
	const op = "[services/restaurants DetailById]"

	restaurant, err := restaurantDTO.ToDB()
	if err != nil {
		log.Warnf("%s Error converting restaurant to DB: %v\n", op, err)
		return &RestaurantReadDTO{}, err
	}

	var newRestaurantId uint
	err = r.db.QueryRowContext(
		ctx,
		"INSERT INTO restaurants (name, link) VALUES ($1, $2) RETURNING id",
		restaurant.Name, restaurant.Link,
	).Scan(&newRestaurantId)
	if err != nil {
		log.Errorf("%s Error inserting restaurant to db: %v\n", op, err)
		return &RestaurantReadDTO{}, err
	}

	newDishDTO := &RestaurantReadDTO{
		ID:   newRestaurantId,
		Name: restaurant.Name,
		Link: restaurant.Link,
	}
	return newDishDTO, nil
}
