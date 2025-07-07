package customers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var customers = []Customer{}

func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(customers); err != nil {
			log.Printf("error encoding customers: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}
	case http.MethodPost:
		var c Customer
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		c.ID = len(customers) + 1
		customers = append(customers, c)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(c); err != nil {
			log.Printf("error encoding customer: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
