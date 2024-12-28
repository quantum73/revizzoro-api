package dto

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type CreateRestaurant struct {
	Name string `json:"name" validate:"required,min=3,max=256"`
	Link string `json:"link" validate:"required,url,max=1024"`
}

func EmptyCreateRestaurant() *CreateRestaurant {
	return &CreateRestaurant{}
}

func (r *CreateRestaurant) GetValue() *CreateRestaurant {
	return r
}

func (r *CreateRestaurant) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s size", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s size", err.Field(), err.Param()))
		case "url":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", err.Field()))
		case "uri":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URI", err.Field()))
		case "uppercase":
			msgs = append(msgs, fmt.Sprintf("%s must be uppercase", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
