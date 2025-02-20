package handlers

import (
	"encoding/json"
	"net/http"
)

type ResBody struct {
	Message string
}

func Index(w http.ResponseWriter, r *http.Request) {

	body := ResBody{Message: "bird-watcher online!"}
	jsonString, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonString)
}
