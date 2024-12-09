package main

import (
	"fmt"
	"github.com/quantum73/revizzoro-api/dishes"
	"github.com/quantum73/revizzoro-api/restaurants"
	"log"
)

func main() {
	log.SetPrefix("[main] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Revizzoro API")

	restaurant, err := restaurants.NewRestaurant(1, "321", "https://321.com")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(restaurant)
	}

	dish, err := dishes.NewDish("Long Bull", 1500, 5, restaurant.Id)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(dish)
	}
}
