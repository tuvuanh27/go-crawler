package service

import (
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
)

type CrawlAliexpressProductPayload struct {
	ProductId string                  `json:"productId"`
	Source    model.ProductTypeSource `json:"source"`
}

type IQueueService interface {
	PublishCrawlAliexpressProduct(payload []CrawlAliexpressProductPayload) error
}

type QueueService struct {
	mqPublisher rabbitmq.IPublisher
	log         logger.ILogger
}

func NewQueueService(mqPublisher rabbitmq.IPublisher) IQueueService {
	return &QueueService{mqPublisher: mqPublisher}
}

func (q *QueueService) PublishCrawlAliexpressProduct(payload []CrawlAliexpressProductPayload) error {
	return q.mqPublisher.PublishMessage(payload)
}
