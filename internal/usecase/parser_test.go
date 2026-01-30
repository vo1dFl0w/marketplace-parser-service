package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/test/mocks"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/usecase"
)

func TestParserService_GetProductsList(t *testing.T) {
	testCases := []struct {
		name      string
		prodName  string
		priceFrom float64
		priceTo   float64
		products  []domain.Product
		expErr    bool
		repoErr   bool
	}{
		{
			name:      "valid",
			prodName:  "prod",
			priceFrom: 0.0,
			priceTo:   500.0,
			products: []domain.Product{
				{
					Name:         "prod",
					Link:         "link1",
					Price:        100.0,
					Rating:       5.0,
					ReviewsCount: 10,
				},
			},
			expErr: false,
		},
		{
			name:      "empty product name",
			prodName:  "",
			priceFrom: 0.0,
			priceTo:   500.0,
			products:  nil,
			expErr:    true,
		},
		{
			name:      "price from below zero",
			prodName:  "prod",
			priceFrom: -100.0,
			priceTo:   500.0,
			products:  nil,
			expErr:    true,
		},
		{
			name:      "price from above price to",
			prodName:  "prod",
			priceFrom: 250.0,
			priceTo:   0.0,
			products:  nil,
			expErr:    true,
		},
		{
			name:      "price to below zero",
			prodName:  "prod",
			priceFrom: 0.0,
			priceTo:   -100.0,
			products:  nil,
			expErr:    true,
		},
	}

	// search service errors
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expErr {
				searchRepo := &mocks.SearchRepositoryMock{}
				searchSrv := usecase.NewSearchService([]repository.SearchRepository{searchRepo})

				searchRepo.On("GetAllProducts", mock.Anything, tc.prodName, tc.priceFrom, tc.priceTo).Return(tc.products, nil)
				res, err := searchSrv.GetProductsList(context.Background(), tc.prodName, tc.priceFrom, tc.priceTo)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.ElementsMatch(t, tc.products, res)

				searchRepo.AssertExpectations(t)
			} else {
				searchRepo := &mocks.SearchRepositoryMock{}
				searchSrv := usecase.NewSearchService([]repository.SearchRepository{searchRepo})

				_, err := searchSrv.GetProductsList(context.Background(), tc.prodName, tc.priceFrom, tc.priceTo)
				assert.Error(t, err)

				searchRepo.AssertNotCalled(t, "GetAllProducts")
			}
		})
	}

	// search repository errors
	t.Run("gateway timeout", func(t *testing.T) {
		p := testCases[0]

		searchRepo := &mocks.SearchRepositoryMock{}
		searchSrv := usecase.NewSearchService([]repository.SearchRepository{searchRepo})

		searchRepo.On("GetAllProducts", mock.Anything, p.prodName, p.priceFrom, p.priceTo).Return(nil, repository.ErrGatewayTimeout).Once()
		_, err := searchSrv.GetProductsList(context.Background(), p.prodName, p.priceFrom, p.priceTo)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrGatewayTimeout)

		searchRepo.AssertExpectations(t)
	})

	t.Run("client closed request", func(t *testing.T) {
		p := testCases[0]

		searchRepo := &mocks.SearchRepositoryMock{}
		searchSrv := usecase.NewSearchService([]repository.SearchRepository{searchRepo})

		searchRepo.On("GetAllProducts", mock.Anything, p.prodName, p.priceFrom, p.priceTo).Return(nil, repository.ErrClientClosedRequest)
		_, err := searchSrv.GetProductsList(context.Background(), p.prodName, p.priceFrom, p.priceTo)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrClientClosedRequest)

		searchRepo.AssertExpectations(t)
	})
}

func TestParserService_ValidateSearchArgs(t *testing.T) {
	testCases := []struct {
		name      string
		prodName  string
		priceFrom float64
		priceTo   float64
		expErr    bool
	}{
		{
			name:      "valid",
			prodName:  "prod",
			priceFrom: 0.0,
			priceTo:   500.0,
			expErr:    false,
		},
		{
			name:      "empty product name",
			prodName:  "",
			priceFrom: 0.0,
			priceTo:   500.0,
			expErr:    true,
		},
		{
			name:      "price from below zero",
			prodName:  "prod",
			priceFrom: -100.0,
			priceTo:   500.0,
			expErr:    true,
		},
		{
			name:      "price from above price to",
			prodName:  "prod",
			priceFrom: 250.0,
			priceTo:   0.0,
			expErr:    true,
		},
		{
			name:      "price to below zero",
			prodName:  "prod",
			priceFrom: 0.0,
			priceTo:   -100.0,
			expErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expErr {
				err := usecase.ValidateSearchArgs(tc.prodName, tc.priceFrom, tc.priceTo)
				assert.NoError(t, err)
			} else {
				err := usecase.ValidateSearchArgs(tc.prodName, tc.priceFrom, tc.priceTo)
				assert.Error(t, err)
			}
		})
	}
}
