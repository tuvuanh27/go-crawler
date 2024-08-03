package configurations

import (
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/consumers/crawl-ali-product"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/delivery"
)

func ConsumerConfigurations(
	mq rabbitmq.IRabbitMQ,
	deliveryBase *delivery.WorkerDeliveryBase,
) error {

	workerConsumer := rabbitmq.NewConsumer[*delivery.WorkerDeliveryBase](deliveryBase.Ctx, deliveryBase.Cfg.Rabbitmq, mq, deliveryBase.Log, deliveryBase.JaegerTracer, crawl_ali_product.HandleCrawlAliexpressProductEvent)

	go func() {
		err := workerConsumer.ConsumeMessage(crawl_ali_product.CrawlAliexpressProductPayload{}, deliveryBase)
		if err != nil {
			deliveryBase.Log.Error(err)
		}
	}()

	return nil
}
