package discounts

import (
	"encoding/json"
	"log"
	"net/http"
)

type Discount struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

var discounts = []Discount{}

func DiscountsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(discounts); err != nil {
			log.Printf("error encoding discounts: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}
	case http.MethodPost:
		var d Discount
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		d.ID = len(discounts) + 1
		discounts = append(discounts, d)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(d); err != nil {
			log.Printf("error encoding discount: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
