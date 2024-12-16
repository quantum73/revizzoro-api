package model

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type Restaurant struct {
	ID   int    `json:"id" validate:"required,gt=0" gorm:"primary_key"`
	Name string `json:"name" validate:"required,max=256" gorm:"not null,size:256"`
	Link string `json:"link" validate:"required,url" gorm:"not null"`
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
