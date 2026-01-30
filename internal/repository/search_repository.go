package repository

import (
	"context"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
)

type SearchRepository interface {
	GetAllProducts(ctx context.Context, name string, priceFrom float64, priceTo float64) ([]domain.Product, error)
}