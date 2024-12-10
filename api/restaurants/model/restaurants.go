package model

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type Restaurant struct {
	Id   int    `json:"id" validate:"required,gt=0"`
	Name string `json:"name" validate:"required,max=256"`
	Link string `json:"link" validate:"required,url"`
}

func (restaurant *Restaurant) GetValue() *Restaurant {
	return restaurant
}

func (restaurant *Restaurant) Validate() error {
	validate := validator.New()
	return validate.Struct(restaurant)
}

func (restaurant *Restaurant) String() string {
	dataAsBytes, err := json.Marshal(&restaurant)
	if err != nil {
		return ""
	}
	return string(dataAsBytes)
}

func NewRestaurant(id int, name, link string) (*Restaurant, error) {
	r := Restaurant{
		Id:   id,
		Name: name,
		Link: link,
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &r, nil
}
