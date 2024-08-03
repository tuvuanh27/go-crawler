package queries

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository/interfaces"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/dtos"
)

type GetProductsHandler struct {
	log               logger.ILogger
	productRepository interfaces.IProductRepository
	ctx               context.Context
}

func NewGetProductsHandler(log logger.ILogger, productRepository interfaces.IProductRepository, ctx context.Context) *GetProductsHandler {
	return &GetProductsHandler{log: log, productRepository: productRepository, ctx: ctx}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProductsQuery) (*dtos.GetPaginationByTypeResponseDto, error) {
	filter := &interfaces.GetProductFilterOptions{
		ProductTypeSource: query.SourceType,
		StartDate:         query.StartDate,
		EndDate:           query.EndDate,
	}

	// Retrieve products from the repository
	if query.Page > 0 {
		products, err := c.productRepository.GetPagination(ctx, query.Page, query.Limit, filter)
		if err != nil {
			return nil, err
		}

		return (*dtos.GetPaginationByTypeResponseDto)(products), nil
	}

	products, err := c.productRepository.GetAllByType(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &dtos.GetPaginationByTypeResponseDto{
		Products: products,
	}, nil
}
