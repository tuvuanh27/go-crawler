package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	echoserver "github.com/tuvuanh27/go-crawler/internal/pkg/http/echo/server"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/services/api/config"
	"go.uber.org/fx"
	"net/http"
)

func RunServers(lc fx.Lifecycle, log logger.ILogger, e *echo.Echo, ctx context.Context, cfg *config.Config, mq rabbitmq.IRabbitMQ) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := echoserver.RunHttpServer(ctx, e, log, cfg.Echo); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running http server: %v", err)
				}
			}()

			//go func() {
			//	if err := grpcServer.RunGrpcServer(ctx); !errors.Is(err, http.ErrServerClosed) {
			//		log.Fatalf("error running grpc server: %v", err)
			//	}
			//}()

			go func() {
				if err := mq.NewRabbitMQConn(ctx, log); err != nil {
					log.Fatalf("error running rabbitmq server: %v", err)
				}
			}()

			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, config.GetMicroserviceName(cfg.ServiceName))
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Infof("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
