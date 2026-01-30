package http_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/test/mocks"
	ht "github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http/httpgen"
)

func TestHandlers_APIV1MarketplaceParserServiceProductsSearchGet(t *testing.T) {
	testCases := []struct {
		name       string
		prodName   string
		priceFrom  float64
		priceTo    float64
		errUsecase error
		expErr     bool
	}{
		{
			name:       "valid",
			prodName:   "prod",
			priceFrom:  0.0,
			priceTo:    500.0,
			errUsecase: nil,
			expErr:     false,
		},
		{
			name:       "empty product name",
			prodName:   "",
			priceFrom:  0.0,
			priceTo:    500.0,
			errUsecase: domain.ErrEmptyProductName,
			expErr:     true,
		},
		{
			name:       "price from below zero",
			prodName:   "prod",
			priceFrom:  -100.0,
			priceTo:    500.0,
			errUsecase: domain.ErrPriceFromBelowZero,
			expErr:     true,
		},
		{
			name:       "price from above price to",
			prodName:   "prod",
			priceFrom:  250.0,
			priceTo:    0.0,
			errUsecase: domain.ErrPriceFromAbovePriceTo,
			expErr:     true,
		},
		{
			name:       "price to below zero",
			prodName:   "prod",
			priceFrom:  0.0,
			priceTo:    -100.0,
			errUsecase: domain.ErrPriceToBelowZero,
			expErr:     true,
		},
		{
			name:       "client closed request",
			prodName:   "prod",
			priceFrom:  0.0,
			priceTo:    500.0,
			errUsecase: domain.ErrClientClosedRequest,
			expErr:     true,
		},
		{
			name:       "gateway timeout",
			prodName:   "prod",
			priceFrom:  0.0,
			priceTo:    500.0,
			errUsecase: domain.ErrGatewayTimeout,
			expErr:     true,
		},
		{
			name:       "internal server error",
			prodName:   "prod",
			priceFrom:  0.0,
			priceTo:    500.0,
			errUsecase: errors.New("internal server error"),
			expErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timeout := time.Second * 30
			parserSrvMock := &mocks.ParserServiceMock{}
			loggerMock := &mocks.LoggerMock{}

			handler := ht.NewHandler(loggerMock, parserSrvMock, timeout)
			if tc.expErr {
				if errors.Is(tc.errUsecase, domain.ErrEmptyProductName) || errors.Is(tc.errUsecase, domain.ErrPriceFromBelowZero) || errors.Is(tc.errUsecase, domain.ErrPriceFromAbovePriceTo) || errors.Is(tc.errUsecase, domain.ErrPriceToBelowZero) {
					parserSrvMock.On("GetProductsList", mock.Anything, tc.prodName, tc.priceFrom, tc.priceTo).Return(nil, tc.errUsecase).Once()
					loggerMock.On("Warn", mock.Anything, mock.Anything).Once()
					res, err := handler.APIV1MarketplaceParserServiceProductsSearchGet(context.Background(), httpgen.APIV1MarketplaceParserServiceProductsSearchGetParams{
						Name:      tc.prodName,
						PriceFrom: httpgen.OptFloat64{Value: tc.priceFrom},
						PriceTo:   httpgen.OptFloat64{Value: tc.priceTo},
					})

					assert.NoError(t, err)
					assert.NotNil(t, res)
					_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetBadRequest)
					assert.True(t, ok)

					parserSrvMock.AssertExpectations(t)
				} else if errors.Is(tc.errUsecase, domain.ErrGatewayTimeout) {
					parserSrvMock.On("GetProductsList", mock.Anything, tc.prodName, tc.priceFrom, tc.priceTo).Return(nil, tc.errUsecase).Once()
					loggerMock.On("Error", mock.Anything, mock.Anything).Once()
					res, err := handler.APIV1MarketplaceParserServiceProductsSearchGet(context.Background(), httpgen.APIV1MarketplaceParserServiceProductsSearchGetParams{
						Name:      tc.prodName,
						PriceFrom: httpgen.OptFloat64{Value: tc.priceFrom},
						PriceTo:   httpgen.OptFloat64{Value: tc.priceTo},
					})

					assert.NoError(t, err)
					assert.NotNil(t, res)
					_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetGatewayTimeout)
					assert.True(t, ok)
				} else if errors.Is(tc.errUsecase, domain.ErrClientClosedRequest) {
					parserSrvMock.On("GetProductsList", mock.Anything, tc.prodName, tc.priceFrom, tc.priceTo).Return(nil, tc.errUsecase).Once()
					loggerMock.On("Warn", mock.Anything, mock.Anything).Once()
					res, err := handler.APIV1MarketplaceParserServiceProductsSearchGet(context.Background(), httpgen.APIV1MarketplaceParserServiceProductsSearchGetParams{
						Name:      tc.prodName,
						PriceFrom: httpgen.OptFloat64{Value: tc.priceFrom},
						PriceTo:   httpgen.OptFloat64{Value: tc.priceTo},
					})

					assert.NoError(t, err)
					assert.NotNil(t, res)
					_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetCode499)
					assert.True(t, ok)

					parserSrvMock.AssertExpectations(t)
				} else {
					parserSrvMock.On("GetProductsList", mock.Anything, tc.prodName, tc.priceFrom, tc.priceTo).Return(nil, tc.errUsecase).Once()
					loggerMock.On("Error", mock.Anything, mock.Anything).Once()
					res, err := handler.APIV1MarketplaceParserServiceProductsSearchGet(context.Background(), httpgen.APIV1MarketplaceParserServiceProductsSearchGetParams{
						Name:      tc.prodName,
						PriceFrom: httpgen.OptFloat64{Value: tc.priceFrom},
						PriceTo:   httpgen.OptFloat64{Value: tc.priceTo},
					})

					assert.NoError(t, err)
					assert.NotNil(t, res)
					_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetInternalServerError)
					assert.True(t, ok)

					parserSrvMock.AssertExpectations(t)
				}
			} else {
				prods := []domain.Product{
					{
						Name:         "prod",
						Link:         "link1",
						Price:        500.0,
						Rating:       5.0,
						ReviewsCount: 253,
					},
				}

				parserSrvMock.On("GetProductsList", mock.Anything, tc.prodName, tc.priceFrom, tc.priceTo).Return(prods, nil).Once()
				res, err := handler.APIV1MarketplaceParserServiceProductsSearchGet(context.Background(), httpgen.APIV1MarketplaceParserServiceProductsSearchGetParams{
					Name:      tc.prodName,
					PriceFrom: httpgen.NewOptFloat64(tc.priceFrom),
					PriceTo:   httpgen.NewOptFloat64(tc.priceTo),
				})

				assert.NoError(t, err)
				assert.NotNil(t, res)

				parserSrvMock.AssertExpectations(t)
			}
		})
	}
}
