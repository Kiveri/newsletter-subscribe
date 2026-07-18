package http_transport

import (
	"context"
	"net/http"
	"time"

	"github.com/Kiveri/newsletter-subscribe/internal/pkg/zap_logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	HeaderRequestID = "X-Request-Id"
)

// Middleware RequestID -> Logger -> RecoverPanic -> Trace -> handler
type Middleware func(handler http.Handler) http.Handler

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(HeaderRequestID)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(HeaderRequestID, requestID)
			w.Header().Set(HeaderRequestID, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *zap_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(HeaderRequestID)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RecoverPanic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := zap_logger.FromCtx(r.Context())
			responseHandler := NewResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := zap_logger.FromCtx(r.Context())
			responseWriter := NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request >>>",
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(responseWriter, r)

			log.Debug(
				"<<< done HTTP request <<<",
				zap.Int("status_code", responseWriter.GetStatusCodeOrPanic()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}
