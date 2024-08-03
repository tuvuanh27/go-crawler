package worker

import (
	"context"
	httpclient "github.com/tuvuanh27/go-crawler/internal/pkg/http-client"
	http "github.com/tuvuanh27/go-crawler/internal/pkg/http/echo"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	mongodriver "github.com/tuvuanh27/go-crawler/internal/pkg/mongo-driver"
	"github.com/tuvuanh27/go-crawler/internal/pkg/otel"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/config"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/configurations"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/delivery"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/service"
	"go.uber.org/fx"
)

func RunRabbitMQ(ctx context.Context, log logger.ILogger, mq rabbitmq.IRabbitMQ) error {
	if err := mq.NewRabbitMQConn(ctx, log); err != nil {
		log.Fatalf("error running rabbitmq server: %v", err)
	}
	return nil
}

func RunWorkers() {
	fx.New(
		fx.Options(
			fx.Provide(
				http.NewContext,
				config.InitConfig,
				logger.InitLogger,
				mongodriver.NewMongo,
				otel.TracerProvider,
				rabbitmq.NewRabbitMQ,
				rabbitmq.NewPublisher,
				httpclient.NewHttpClient,
				service.NewOpenAliexpressService,
				repository.NewProductRepository,
				delivery.RunDelivery,
			),
			fx.Invoke(RunRabbitMQ),
			fx.Invoke(configurations.ConsumerConfigurations),
		),
	).Run()
}
