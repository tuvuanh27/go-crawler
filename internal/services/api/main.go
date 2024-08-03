package main

import (
	"github.com/go-playground/validator"
	http "github.com/tuvuanh27/go-crawler/internal/pkg/http/echo"
	echoserver "github.com/tuvuanh27/go-crawler/internal/pkg/http/echo/server"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	mongodriver "github.com/tuvuanh27/go-crawler/internal/pkg/mongo-driver"
	"github.com/tuvuanh27/go-crawler/internal/pkg/otel"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository"
	"github.com/tuvuanh27/go-crawler/internal/services/api/config"
	"github.com/tuvuanh27/go-crawler/internal/services/api/configurations"
	"github.com/tuvuanh27/go-crawler/internal/services/api/mappings"
	"github.com/tuvuanh27/go-crawler/internal/services/api/server"
	"github.com/tuvuanh27/go-crawler/internal/services/api/service"
	"go.uber.org/fx"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				http.NewContext,
				config.InitConfig,
				logger.InitLogger,
				echoserver.NewEchoServer,
				//grpc.NewGrpcServer,
				mongodriver.NewMongo,
				otel.TracerProvider,
				rabbitmq.NewRabbitMQ,
				rabbitmq.NewPublisher,
				validator.New,
				repository.NewUserRepository,
				repository.NewProductRepository,
				service.NewQueueService,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(configurations.ConfigMiddlewares),
			fx.Invoke(configurations.ConfigSwagger),
			fx.Invoke(mappings.ConfigureMappings),
			fx.Invoke(configurations.ConfigEndpoints),
			fx.Invoke(configurations.ConfigUsersMediator),
			fx.Invoke(configurations.ConfigProductMediator),
		),
	).Run()
}
