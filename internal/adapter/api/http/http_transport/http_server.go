package http_transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Kiveri/newsletter-subscribe/internal/config/http_config"
	"github.com/Kiveri/newsletter-subscribe/internal/pkg/zap_logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config http_config.Config
	log    *zap_logger.Logger
}

func NewHTTPServer(config http_config.Config, log *zap_logger.Logger) *HTTPServer {
	return &HTTPServer{
		mux:    http.NewServeMux(),
		config: config,
		log:    log,
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: h.mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Warn("starting HTTP server", zap.String("addr", h.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		h.log.Warn("HTTP server shutdown completed")
	}

	return nil
}

func (h *HTTPServer) RegisterApiRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := "/api" + string(router.apiVersion)

		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}
