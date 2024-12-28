package dto

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type CreateDish struct {
	Name         string `json:"name" validate:"required,min=3,max=256"`
	Price        int    `json:"price" validate:"required,gt=0"`
	Score        int    `json:"score" validate:"required,gt=0,lte=5"`
	RestaurantID int    `json:"restaurant_id" validate:"required,gt=0"`
}

func EmptyCreateDish() *CreateDish {
	return &CreateDish{}
}

func (r *CreateDish) GetValue() *CreateDish {
	return r
}

func (r *CreateDish) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s size", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s size", err.Field(), err.Param()))
		case "gt":
			msgs = append(msgs, fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param()))
		case "lte":
			msgs = append(msgs, fmt.Sprintf("%s must be less or equal %s", err.Field(), err.Param()))
		case "url":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
