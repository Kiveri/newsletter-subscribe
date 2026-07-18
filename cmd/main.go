package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kiveri/newsletter-subscribe/internal/adapter/api/http/http_transport"
	"github.com/Kiveri/newsletter-subscribe/internal/adapter/api/http/user_http_handler"
	"github.com/Kiveri/newsletter-subscribe/internal/config/http_config"
	"github.com/Kiveri/newsletter-subscribe/internal/pkg/zap_logger"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := zap_logger.NewLogger(zap_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app zap_logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Starting app!")

	usersHandler := user_http_handler.NewHandler(nil)
	usersRoutes := usersHandler.Routes()

	apiVersionRouter := http_transport.NewApiVersionRouter(http_transport.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := http_transport.NewHTTPServer(
		http_config.NewConfigMust(),
		logger,
	)
	httpServer.RegisterApiRouters(apiVersionRouter)

	if err = httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
