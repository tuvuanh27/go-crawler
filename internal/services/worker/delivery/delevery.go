package delivery

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository/interfaces"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/config"
	service "github.com/tuvuanh27/go-crawler/internal/services/worker/service/interface"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

type WorkerDeliveryBase struct {
	Log               logger.ILogger
	Cfg               *config.Config
	RabbitmqPublisher rabbitmq.IPublisher
	JaegerTracer      trace.Tracer
	Mongo             *mongo.Database
	Ctx               context.Context
	AliexpressService service.AliexpressService
	ProductRepository interfaces.IProductRepository
}

func RunDelivery(ctx context.Context, cfg *config.Config, log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher, mongo *mongo.Database, jaegerTracer trace.Tracer, aliexpressService service.AliexpressService, productRepository interfaces.IProductRepository) *WorkerDeliveryBase {
	return &WorkerDeliveryBase{
		Log:               log,
		Cfg:               cfg,
		RabbitmqPublisher: rabbitmqPublisher,
		JaegerTracer:      jaegerTracer,
		Mongo:             mongo,
		Ctx:               ctx,
		AliexpressService: aliexpressService,
		ProductRepository: productRepository,
	}
}
