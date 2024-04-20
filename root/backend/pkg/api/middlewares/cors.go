package middlewares

import "net/http"

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	  w.Header().Set("Access-Control-Allow-Origin", "*")
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
