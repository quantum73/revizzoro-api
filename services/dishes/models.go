package dishes

import "encoding/json"

type dishDB struct {
	ID           uint
	Name         string
	Price        uint
	Score        uint
	RestaurantID uint
}

func (d dishDB) toDTO() *DishDTO {
	return &DishDTO{
		ID:           d.ID,
		Name:         d.Name,
		Price:        d.Price,
		Score:        d.Score,
		RestaurantID: d.RestaurantID,
	}
}

type DishDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Price        uint   `json:"price"`
	Score        uint   `json:"score"`
	RestaurantID uint   `json:"restaurant_id"`
}

func (dto *DishDTO) MarshallJSON() ([]byte, error) {
	return json.Marshal(dto)
}
