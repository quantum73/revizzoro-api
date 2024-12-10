package dishes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quantum73/revizzoro-api/api/dishes/model"
	"net/http"
	"strconv"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Dishes list")
}

func DetailByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "Dish not found")
	}

	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, "Dish not found")
	}

	newDish, err := model.NewDish(idAsInt, "Long Bull", 1500, 5, 1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, newDish.String())
}
