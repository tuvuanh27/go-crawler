package commands

import "github.com/tuvuanh27/go-crawler/internal/pkg/model"

type CrawlAliexpressProduct struct {
	ProductIds []string                `json:"productIds" required:"true" validate:"required"`
	Source     model.ProductTypeSource `json:"source" required:"true" validate:"required"`
}

func NewCrawlAliexpressProduct(productIds []string, source model.ProductTypeSource) *CrawlAliexpressProduct {
	return &CrawlAliexpressProduct{
		ProductIds: productIds,
		Source:     source,
	}
}
