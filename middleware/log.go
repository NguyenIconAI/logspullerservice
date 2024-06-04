package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// LogMiddleware logs the request and response details
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Initialize the status to 200 in case WriteHeader is not called
		rec := statusRecorder{w, 200}

		id := time.Now().UnixNano()
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		// log request details
		log.Printf("%d %s %s", id, method, uri)

		next.ServeHTTP(&rec, r)

		duration := time.Since(start)
		status := rec.status

		// log response details
		log.Printf("%d %d %s", id, status, duration)

	})
}
