package locations

import (
	"encoding/json"
	"log"
	"net/http"
)

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var locations = []Location{}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(locations); err != nil {
			log.Printf("error encoding locations: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}
	case http.MethodPost:
		var l Location
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		l.ID = len(locations) + 1
		locations = append(locations, l)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(l); err != nil {
			log.Printf("error encoding location: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
