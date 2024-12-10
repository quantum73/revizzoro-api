package restaurants

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quantum73/revizzoro-api/api/restaurants/model"
	"net/http"
	"strconv"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Restaurants list")
}

func DetailByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "Restaurant not found")
	}

	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "Restaurant not found")
	}

	newRestaurant, err := model.NewRestaurant(idAsInt, "MockName", "http://some-rest.com")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, newRestaurant.String())
}
