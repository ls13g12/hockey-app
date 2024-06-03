package api


import (
	"log/slog"
	"net/http"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func (a *api) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode: http.StatusOK,
		}
		next.ServeHTTP(rec, r)
		duration := time.Since(startTime)

		if rec.StatusCode > 299 {
			a.logger.Error(
				"request",
				slog.String("method", r.Method),
				slog.String("path", r.RequestURI),
				slog.Int("status_code", rec.StatusCode),
				slog.String("body", string(rec.Body[:])),
				slog.String("status_text", http.StatusText(rec.StatusCode)),
				slog.Duration("duration", duration),
			)
			return
		}

		a.logger.Info(
			"request",
			slog.String("method", r.Method),
			slog.String("path", r.RequestURI),
			slog.Int("status_code", rec.StatusCode),
			slog.String("status_text", http.StatusText(rec.StatusCode)),
			slog.Duration("duration", duration),
		)
	})
}

func (a *api) corsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  w.Header().Set("Access-Control-Allow-Origin", "*")
	  w.Header().Add("Access-Control-Allow-Credentials", "true")
	  w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	  w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  
	  if r.Method == "OPTIONS" {
		  http.Error(w, "No Content", http.StatusNoContent)
		  return
	  }
  
	  next.ServeHTTP(w, r)
	})
}

func (a *api) isAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !a.sessionManager.Exists(r.Context(), "authenticatedUserID") {
			http.Error(w, "Not Authorised", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
