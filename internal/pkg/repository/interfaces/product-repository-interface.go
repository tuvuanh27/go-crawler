package interfaces

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
)

type GetPaginationByTypeResponse struct {
	Products []*model.Product `json:"products"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

type GetProductFilterOptions struct {
	ProductTypeSource model.ProductTypeSource
	StartDate         uint
	EndDate           uint
}

type IProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*model.Product, error)
	CreateMany(ctx context.Context, products []*model.Product) ([]*model.Product, error)
	GetPagination(ctx context.Context, page, limit uint, filterOptions *GetProductFilterOptions) (*GetPaginationByTypeResponse, error)
	GetAllByType(ctx context.Context, filterOptions *GetProductFilterOptions) ([]*model.Product, error)
}
