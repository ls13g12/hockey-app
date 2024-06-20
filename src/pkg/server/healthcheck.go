package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *server) healthcheckGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("error is %v", err)
	}
	w.Write(jsonData)
}
