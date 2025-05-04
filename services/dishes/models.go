package dishes

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type DishDB struct {
	ID           uint
	Name         string
	Price        uint
	Score        uint
	RestaurantID uint
}

func (d DishDB) ToReadDTO() (*DishReadDTO, error) {
	dto := &DishReadDTO{
		ID:           d.ID,
		Name:         d.Name,
		Price:        d.Price,
		Score:        d.Score,
		RestaurantID: d.RestaurantID,
	}
	err := dto.Validate()
	if err != nil {
		return &DishReadDTO{}, err
	}
	return dto, nil
}

func (d DishDB) ToCreateDTO() (*DishCreateDTO, error) {
	dto := &DishCreateDTO{
		Name:         d.Name,
		Price:        d.Price,
		Score:        d.Score,
		RestaurantID: d.RestaurantID,
	}
	err := dto.Validate()
	if err != nil {
		return &DishCreateDTO{}, err
	}
	return dto, nil
}

type DishCreateDTO struct {
	Name         string `json:"name" validate:"required"`
	Price        uint   `json:"price" validate:"required,gt=0"`
	Score        uint   `json:"score" validate:"required,gt=0,lte=5"`
	RestaurantID uint   `json:"restaurant_id" validate:"required,gt=0"`
}

func (d DishCreateDTO) Validate() error {
	return validate.Struct(d)
}

func (d DishCreateDTO) ToDB() (*DishDB, error) {
	err := validate.Struct(d)
	if err != nil {
		return &DishDB{}, err
	}

	dish := &DishDB{
		Name:         d.Name,
		Price:        d.Price,
		Score:        d.Score,
		RestaurantID: d.RestaurantID,
	}
	return dish, nil
}

type DishReadDTO struct {
	ID           uint   `json:"id" validate:"required,gt=0"`
	Name         string `json:"name" validate:"required"`
	Price        uint   `json:"price" validate:"required,gt=0"`
	Score        uint   `json:"score" validate:"required,gt=0,lte=5"`
	RestaurantID uint   `json:"restaurant_id" validate:"required,gt=0"`
}

func (d DishReadDTO) Validate() error {
	return validate.Struct(d)
}
