package crawl_ali_product

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/delivery"
	"time"
)

func HandleCrawlAliexpressProductEvent(queue string, msg amqp.Delivery, deliveryBase *delivery.WorkerDeliveryBase) error {
	// log the message body
	deliveryBase.Log.Info("Received message: ", string(msg.Body))

	// parse the message body to CrawlAliexpressProductPayload
	var payload []CrawlAliexpressProductPayload
	err := json.Unmarshal(msg.Body, &payload)
	if err != nil {
		deliveryBase.Log.Error(err)
		return err
	}

	deliveryBase.Log.Info("Parse message: ", payload)
	var products []*model.Product

	for _, v := range payload {
		product, err := deliveryBase.AliexpressService.GetProduct(v.ProductId)
		if err != nil {
			deliveryBase.Log.Error(err)
			return err
		}
		// append product to products
		products = append(products, product)
		time.Sleep(100 * time.Millisecond)
	}

	//// save products to database
	_, err = deliveryBase.ProductRepository.CreateMany(deliveryBase.Ctx, products)
	if err != nil {
		deliveryBase.Log.Error(err)
		return err
	}

	return nil
}
