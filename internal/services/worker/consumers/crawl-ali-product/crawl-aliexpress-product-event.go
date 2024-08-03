package crawl_ali_product

import "github.com/tuvuanh27/go-crawler/internal/pkg/model"

type CrawlAliexpressProductPayload struct {
	ProductId string                  `json:"productId"`
	Source    model.ProductTypeSource `json:"source"`
}
