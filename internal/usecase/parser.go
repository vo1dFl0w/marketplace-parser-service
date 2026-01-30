package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
)

type ParserService interface {
	GetProductsList(ctx context.Context, name string, priceFrom float64, priceTo float64) ([]domain.Product, error)
}

type parserService struct {
	source []repository.SearchRepository
}

func NewSearchService(source []repository.SearchRepository) *parserService {
	return &parserService{source: source}
}

func (s *parserService) GetProductsList(ctx context.Context, name string, priceFrom float64, priceTo float64) ([]domain.Product, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := ValidateSearchArgs(name, priceFrom, priceTo); err != nil {
		return nil, err
	}

	resCh := make(chan []domain.Product)
	errCh := make(chan error, 1)

	wg := &sync.WaitGroup{}
	for _, src := range s.source {
		wg.Add(1)
		go func(source repository.SearchRepository) {
			defer wg.Done()
			products, err := source.GetAllProducts(ctx, name, priceFrom, priceTo)
			if err != nil {
				if errors.Is(err, repository.ErrGatewayTimeout) {
					select {
					case errCh <- domain.ErrGatewayTimeout:
						break
					default:
					}
					cancel()
					return
				} else if errors.Is(err, repository.ErrClientClosedRequest) {
					select {
					case errCh <- domain.ErrClientClosedRequest:
						break
					default:
					}
					cancel()
					return
				} else {
					select {
					case errCh <- err:
						break
					default:
					}
					cancel()
					return
				}
			}

			select {
			case resCh <- products:
				return
			case <-ctx.Done():
				return
			}
		}(src)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	res := []domain.Product{}
	for {
		select {
		case err := <-errCh:
			return nil, err
		case r, ok := <-resCh:
			if !ok {
				return res, nil
			}
			res = append(res, r...)
		case <-ctx.Done():
			if ctx.Err() != nil {
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return nil, domain.ErrGatewayTimeout
				}
				if errors.Is(ctx.Err(), context.Canceled) {
					return nil, domain.ErrClientClosedRequest
				}
				return nil, ctx.Err()
			}
		}
	}
}

func ValidateSearchArgs(name string, priceFrom float64, priceTo float64) error {
	if name == "" {
		return domain.ErrEmptyProductName
	}

	if priceFrom < 0 {
		return domain.ErrPriceFromBelowZero
	}

	if priceFrom > priceTo {
		return domain.ErrPriceFromAbovePriceTo
	}

	if priceTo < 0 {
		return domain.ErrPriceToBelowZero
	}

	return nil
}
