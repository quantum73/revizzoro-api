package restaurants

import "encoding/json"

type Restaurant struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

func (r *Restaurant) String() string {
	dataAsBytes, err := json.MarshalIndent(&r, "", "  ")
	if err != nil {
		return ""
	}
	return string(dataAsBytes)
}

func NewRestaurant(id int, name, link string) (*Restaurant, error) {
	err := validateId(id)
	if err != nil {
		return nil, err
	}

	err = validateName(name)
	if err != nil {
		return nil, err
	}

	return &Restaurant{Id: id, Name: name, Link: link}, nil
}
