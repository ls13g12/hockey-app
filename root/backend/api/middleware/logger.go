package middleware

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

func HttpLogger(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode: http.StatusOK,
		}
		next.ServeHTTP(rec, r)
		duration := time.Since(startTime)

		if rec.StatusCode != http.StatusOK {
			logger.Error(
				"Http request",
				slog.String("method", r.Method),
				slog.String("path", r.RequestURI),
				slog.Int("status_code", rec.StatusCode),
				slog.String("body", string(rec.Body[:])),
				slog.String("status_text", http.StatusText(rec.StatusCode)),
				slog.Duration("duration", duration),
			)
			return
		}

		logger.Info(
			"Http request",
			slog.String("method", r.Method),
			slog.String("path", r.RequestURI),
			slog.Int("status_code", rec.StatusCode),
			slog.String("status_text", http.StatusText(rec.StatusCode)),
			slog.Duration("duration", duration),
		)
	})
}
