package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (h *Handler) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := h.logger.With(
			"remote_addr", r.RemoteAddr,
			"http-method", r.Method,
			"path", r.URL.Path,
		)

		log.Info("started")

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		completed := time.Since(start)
		completedStr := fmt.Sprintf("%.3fms", float64(completed.Microseconds())/1000)

		statusText := http.StatusText(rw.code)
		if rw.code == 499 {
			statusText = ErrClientClosedRequest.Error()
		}
		attrs := []any{
			"code", rw.code,
			"status-text", statusText,
			"duration_ms", completedStr,
		}

		switch {
		case rw.code >= 500:
			log.Error("failed", attrs...)
		case rw.code >= 400:
			log.Warn("failed", attrs...)
		default:
			log.Info("completed", attrs...)
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (h *Handler) RequestTimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.requestTimeout)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
