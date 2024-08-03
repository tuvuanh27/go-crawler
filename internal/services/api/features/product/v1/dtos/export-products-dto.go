package dtos

import (
	"github.com/pkg/errors"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
)

type ExportType string

var (
	ExportFull ExportType = "full"
	ExportMain ExportType = "main"
)

type ExportAliexpressProductsRequestDto struct {
	SourceType uint       `query:"type" required:"true" validate:"required"`
	StartDate  uint       `query:"start_date"`
	EndDate    uint       `query:"end_date"`
	ExportType ExportType `query:"export_type" validate:"required"`
}

func (r *ExportAliexpressProductsRequestDto) Validate() error {
	if _, ok := model.ValidProductTypeSources[r.SourceType]; !ok {
		return errors.New("invalid source type")
	}

	// validate start date and end date
	if r.StartDate > 0 && r.EndDate > 0 && r.StartDate > r.EndDate {
		return errors.New("Start date must be less than end date")
	}

	if r.ExportType != ExportFull && r.ExportType != ExportMain {
		return errors.New("Invalid export type")
	}

	return nil
}
