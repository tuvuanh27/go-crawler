package service

import "github.com/tuvuanh27/go-crawler/internal/pkg/model"

type AliexpressService interface {
	GetProduct(productID string) (*model.Product, error)
}
