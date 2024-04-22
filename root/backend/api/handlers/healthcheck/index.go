package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
)


func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /healthcheck request\n")
	w.Header().Set("Content-Type", "application/json")
	payload := struct{
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
