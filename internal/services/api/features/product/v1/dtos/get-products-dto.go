package dtos

import (
	"github.com/pkg/errors"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
)

type GetProductsRequestDto struct {
	SourceType uint `query:"type" required:"true" validate:"required"`
	StartDate  uint `query:"start_date"`
	EndDate    uint `query:"end_date"`
	Page       uint `query:"page"`
	Limit      uint `query:"limit"`
}

func (r *GetProductsRequestDto) Validate() error {
	if _, ok := model.ValidProductTypeSources[r.SourceType]; !ok {
		return errors.New("invalid source type")
	}

	// validate start date and end date
	if r.StartDate > 0 && r.EndDate > 0 && r.StartDate > r.EndDate {
		return errors.New("start date must be less than end date")
	}

	if r.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	return nil
}

type GetPaginationByTypeResponseDto struct {
	Products []*model.Product `json:"products"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}
