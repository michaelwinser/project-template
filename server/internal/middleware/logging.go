package middleware

import (
	"net/http"
	"time"

	"project-template/server/internal/logging"
)

// responseWriter wraps http.ResponseWriter to capture status code and size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Logging logs HTTP requests
func Logging(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := GetRequestID(r.Context())

			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			ctx := map[string]interface{}{
				"request_id":    requestID,
				"method":        r.Method,
				"path":          r.URL.Path,
				"status":        wrapped.statusCode,
				"duration_ms":   duration.Milliseconds(),
				"response_size": wrapped.size,
				"remote":        r.RemoteAddr,
			}

			// Add user agent if present
			if ua := r.Header.Get("User-Agent"); ua != "" {
				ctx["user_agent"] = ua
			}

			// Add request content length for POST/PUT
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				ctx["request_size"] = r.ContentLength
			}

			logger.Info("request completed", ctx)
		})
	}
}
