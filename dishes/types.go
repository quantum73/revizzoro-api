package dishes

import (
	"encoding/json"
	"fmt"
)

type Dish struct {
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Score        int    `json:"score"`
	RestaurantId int    `json:"restaurant_id"`
}

func (d *Dish) String() string {
	dataAsBytes, err := json.MarshalIndent(&d, "", "  ")
	if err != nil {
		return fmt.Sprintf("Dish[%s]", d.Name)
	}
	return string(dataAsBytes)
}

func NewDish(name string, price, score, restaurantId int) (*Dish, error) {
	err := validateName(name)
	if err != nil {
		return nil, err
	}
	err = validatePrice(price)
	if err != nil {
		return nil, err
	}
	err = validateScore(score)
	if err != nil {
		return nil, err
	}

	return &Dish{
		Name:         name,
		Price:        price,
		Score:        score,
		RestaurantId: restaurantId,
	}, nil
}
