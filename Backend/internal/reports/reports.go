package reports

import (
	"encoding/json"
	"log"
	"net/http"
)

type Report struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Data  string `json:"data"`
}

var reports = []Report{}

func ReportsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(reports); err != nil {
			log.Printf("error encoding reports: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}
	case http.MethodPost:
		var rep Report
		if err := json.NewDecoder(r.Body).Decode(&rep); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid input"}`))
			return
		}
		rep.ID = len(reports) + 1
		reports = append(reports, rep)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(rep); err != nil {
			log.Printf("error encoding report: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
