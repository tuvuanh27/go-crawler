package queries

import "github.com/tuvuanh27/go-crawler/internal/pkg/model"

type GetProductsQuery struct {
	SourceType model.ProductTypeSource `json:"sourceType" required:"true" validate:"required"`
	StartDate  uint                    `json:"startDate"`
	EndDate    uint                    `json:"endDate"`
	Page       uint                    `json:"page"`
	Limit      uint                    `json:"limit"`
}

func NewGetProductsQuery(query GetProductsQuery) *GetProductsQuery {
	return &GetProductsQuery{
		SourceType: query.SourceType,
		StartDate:  query.StartDate,
		EndDate:    query.EndDate,
		Page:       query.Page,
		Limit:      query.Limit,
	}
}
