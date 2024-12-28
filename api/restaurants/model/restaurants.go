package model

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	dishes "github.com/quantum73/revizzoro-api/api/dishes/model"
)

type Restaurant struct {
	ID     int           `json:"id,omitempty" gorm:"primary_key,auto_increment"`
	Name   string        `json:"name" validate:"required,min=3,max=256" gorm:"not null,size:256"`
	Link   string        `json:"link" validate:"required,url,max=1024" gorm:"not null,size:1024"`
	Dishes []dishes.Dish `json:"dishes,omitempty" gorm:"foreignKey:RestaurantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func NewRestaurant(name, link string) (*Restaurant, error) {
	r := Restaurant{Name: name, Link: link}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &r, nil
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
