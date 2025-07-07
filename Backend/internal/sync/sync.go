package sync

import (
	"encoding/json"
	"log"
	"net/http"
)

type SyncStatus struct {
	Status string `json:"status"`
}

func SyncHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		status := SyncStatus{Status: "ok (placeholder)"}
		if err := json.NewEncoder(w).Encode(status); err != nil {
			log.Printf("error encoding sync status: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal error"}`))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
