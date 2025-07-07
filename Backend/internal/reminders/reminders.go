package reminders

import (
	"encoding/json"
	"log"
	"net/http"
)

type Reminder struct {
	ID      int    `json:"id"`
	Booking int    `json:"booking"`
	Message string `json:"message"`
}

var reminders = []Reminder{}

func RemindersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(reminders); err != nil {
			log.Printf("error encoding reminders: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}
	case http.MethodPost:
		var rem Reminder
		if err := json.NewDecoder(r.Body).Decode(&rem); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		rem.ID = len(reminders) + 1
		reminders = append(reminders, rem)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(rem); err != nil {
			log.Printf("error encoding reminder: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
