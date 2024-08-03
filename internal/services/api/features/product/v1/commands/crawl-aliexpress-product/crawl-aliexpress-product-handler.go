package commands

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/services/api/service"
)

type CrawlAliexpressProductHandler struct {
	log          logger.ILogger
	ctx          context.Context
	queueService service.IQueueService
}

func NewCrawlAliexpressProductHandler(log logger.ILogger, ctx context.Context, queueService service.IQueueService) *CrawlAliexpressProductHandler {
	return &CrawlAliexpressProductHandler{log: log, ctx: ctx, queueService: queueService}
}

func (c *CrawlAliexpressProductHandler) Handle(ctx context.Context, command *CrawlAliexpressProduct) (*string, error) {
	msg := "Add crawl products successfully"
	productIds := command.ProductIds
	source := command.Source

	var payload []service.CrawlAliexpressProductPayload
	for _, v := range productIds {
		payload = append(payload, service.CrawlAliexpressProductPayload{ProductId: v, Source: source})
	}

	// chunk payload to 10 items
	if len(payload) > 20 {
		for i := 0; i < len(payload); i += 20 {
			end := i + 20
			if end > len(payload) {
				end = len(payload)
			}
			if err := c.queueService.PublishCrawlAliexpressProduct(payload[i:end]); err != nil {
				return nil, err
			}
		}
		return &msg, nil
	}

	if err := c.queueService.PublishCrawlAliexpressProduct(payload); err != nil {
		return nil, err
	}

	return &msg, nil
}
