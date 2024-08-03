package dtos

import "github.com/tuvuanh27/go-crawler/internal/pkg/model"

type CrawlDto struct {
	ProductIds []string                `json:"productIds" required:"true" validate:"required"`
	Source     model.ProductTypeSource `json:"source" required:"true" validate:"required"`
}
