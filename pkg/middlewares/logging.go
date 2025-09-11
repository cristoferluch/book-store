package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		zap.L().Info("http request",
			zap.String("method", r.Method),
			zap.String("host", r.Host),
			zap.String("url", r.URL.String()),
			zap.Int("status", ww.Status()),
			zap.String("latency", time.Since(start).String()),
			zap.String("ip", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}
