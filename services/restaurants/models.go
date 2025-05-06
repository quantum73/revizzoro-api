package restaurants

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type RestaurantDB struct {
	ID   uint
	Name string
	Link string
}

func (d RestaurantDB) ToReadDTO() (*RestaurantReadDTO, error) {
	dto := &RestaurantReadDTO{
		ID:   d.ID,
		Name: d.Name,
		Link: d.Link,
	}
	err := dto.Validate()
	if err != nil {
		return &RestaurantReadDTO{}, err
	}
	return dto, nil
}

func (d RestaurantDB) ToCreateDTO() (*RestaurantCreateDTO, error) {
	dto := &RestaurantCreateDTO{
		Name: d.Name,
		Link: d.Link,
	}
	err := dto.Validate()
	if err != nil {
		return &RestaurantCreateDTO{}, err
	}
	return dto, nil
}

type RestaurantCreateDTO struct {
	Name string `json:"name" validate:"required,min=3,max=256"`
	Link string `json:"link" validate:"required,url,min=3,max=1024"`
}

func (d RestaurantCreateDTO) Validate() error {
	return validate.Struct(d)
}

func (d RestaurantCreateDTO) ToDB() (*RestaurantDB, error) {
	err := validate.Struct(d)
	if err != nil {
		return &RestaurantDB{}, err
	}

	dish := &RestaurantDB{
		Name: d.Name,
		Link: d.Link,
	}
	return dish, nil
}

type RestaurantReadDTO struct {
	ID   uint   `json:"id" validate:"required,gt=0"`
	Name string `json:"name" validate:"required,min=3,max=256"`
	Link string `json:"link" validate:"required,url,min=3,max=1024"`
}

func (d RestaurantReadDTO) Validate() error {
	return validate.Struct(d)
}
