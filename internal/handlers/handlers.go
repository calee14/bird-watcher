package handlers

import (
	db "bird-watcher/internal/database"
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	body := Response{Message: "bird-watcher online!"}
	jsonString, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonString)
}

func AddSubscriber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subscriber db.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subscriber); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubscriber := db.NewSubscriber(subscriber.Email)
	if err := db.CreateSubscriber(newSubscriber); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSubscriber)
}

func RemoveSusbcriber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subscriber db.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&subscriber); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.DeleteSubscriber(subscriber.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Message: "successfully unsubscribed"})
}
