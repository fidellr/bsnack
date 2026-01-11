package middleware

import (
	"bsnack/pkg/logger"
	"net/http"
	"time"
)

// responseWriter is a wrapper to capture the HTTP status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		logger.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrappedWriter.status,
			"duration", time.Since(start).String(),
			"remote_ip", r.RemoteAddr,
		)
	})
}
