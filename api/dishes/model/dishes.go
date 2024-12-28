package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Dish struct {
	ID           int    `json:"id,omitempty" gorm:"primary_key,auto_increment"`
	Name         string `json:"name" validate:"required,min=3,max=256" gorm:"not null,size:256"`
	Price        int    `json:"price" validate:"required,gt=0" gorm:"not null,check:name > 0"`
	Score        int    `json:"score" validate:"required,gt=0,lte=5" gorm:"not null,check:score > 0 && score <= 5"`
	RestaurantID int    `json:"restaurant_id" validate:"required,gt=0" gorm:"not null"`
}

func NewDish(name string, price, score, restaurantId int) (*Dish, error) {
	d := Dish{
		Name:         name,
		Price:        price,
		Score:        score,
		RestaurantID: restaurantId,
	}
	if err := d.Validate(); err != nil {
		return nil, err
	}
	return &d, nil
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
