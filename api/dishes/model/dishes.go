package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Dish struct {
	Id           int    `json:"id" validate:"required,gt=0"`
	Name         string `json:"name" validate:"required,max=256"`
	Price        int    `json:"price" validate:"required,gt=0"`
	Score        int    `json:"score" validate:"required,gt=0,lte=5"`
	RestaurantId int    `json:"restaurant_id" validate:"required,gt=0"`
}

func (dish *Dish) GetValue() *Dish {
	return dish
}

func (dish *Dish) Validate() error {
	validate := validator.New()
	return validate.Struct(dish)
}

func (dish *Dish) String() string {
	dataAsBytes, err := json.Marshal(&dish)
	if err != nil {
		return fmt.Sprintf("Dish[%s]", dish.Name)
	}
	return string(dataAsBytes)
}

func NewDish(id int, name string, price, score, restaurantId int) (*Dish, error) {
	d := Dish{
		Id:           id,
		Name:         name,
		Price:        price,
		Score:        score,
		RestaurantId: restaurantId,
	}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	return &d, nil
}
