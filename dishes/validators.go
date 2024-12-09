package dishes

import (
	"errors"
	"fmt"
)

const (
	minScore = 1
	maxScore = 5
)

func validateName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}

func validatePrice(price int) error {
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}

func validateScore(score int) error {
	if score < minScore || score > maxScore {
		return fmt.Errorf(
			"score must be equal or bigger than %d and less or equal %d",
			minScore, maxScore,
		)
	}
	return nil
}
