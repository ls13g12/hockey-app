package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type getHelloResp struct {
	WelcomeMsg string `json:"welcome_msg"`
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Response from /"))
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	w.Header().Set("Content-Type", "application/json")
	data := getHelloResp{
		WelcomeMsg: "Hello Hockey App",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error is %v", err)
	}
	w.Write(jsonData)
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	  w.Header().Add("Access-Control-Allow-Origin", "*")
	  w.Header().Add("Access-Control-Allow-Credentials", "true")
	  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	  w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  
	  if r.Method == "OPTIONS" {
		  http.Error(w, "No Content", http.StatusNoContent)
		  return
	  }
  
	  next(w, r)
	}
  }

func main() {
	http.HandleFunc("/", CORS(getRoot))
	http.HandleFunc("/hello", CORS(getHello))

	fmt.Print("Starting server at port 3333\n")
	http.ListenAndServe(":3333", nil)
}