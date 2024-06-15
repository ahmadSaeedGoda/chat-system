package handlers

import (
	"encoding/json"
	"net/http"
)

func HowAmI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := struct {
		Status, HeartBeat, Message string
	}{
		Status:    "ALIVE",
		HeartBeat: "HeartBeating & Functioning",
		Message:   "How Are you?!",
	}

	json.NewEncoder(w).Encode(res)
}
