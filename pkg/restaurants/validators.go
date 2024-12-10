package restaurants

import (
	"errors"
)

func validateId(id int) error {
	if id <= 0 {
		return errors.New("id must be a integer bigger than 0")
	}
	return nil
}

func validateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}
