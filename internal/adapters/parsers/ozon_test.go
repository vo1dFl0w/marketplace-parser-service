package parsers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/parsers"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/test/mocks"
)

func TestParsers_NewOzonParser(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			OzonCfg: &config.OzonConfig{
				BaseURL:           "url",
				SearchBarSelector: "search-ber-selector",
				ItemsSelector:     "items-selector",
				LinkSelector:      "link-selector",
				PriceSelector:     "price-selector",
				RatingSelector:    "rating-selector",
				ReviewsSelector:   "reviews-selector",
			},
		},
	}

	browserRepoMock := &mocks.BrowserRepositoryMock{}

	ozonParser := parsers.NewOzonParser(cfg, nil, browserRepoMock)
	assert.NotNil(t, ozonParser)
}

func TestParsers_OzonParser(t *testing.T) {
	loggerMock := &mocks.LoggerMock{}
	browserRepoMock := &mocks.BrowserRepositoryMock{}
	pageMock := &mocks.PageMock{}

	searchBarMock := &mocks.ElementMock{}
	itemMock := &mocks.ElementMock{}
	linkElMock := &mocks.ElementMock{}
	priceElMock := &mocks.ElementMock{}
	prodNameElMock := &mocks.ElementMock{}
	ratingElMock := &mocks.ElementMock{}
	reviewsElMock := &mocks.ElementMock{}

	cfg := &config.Config{
		Server: config.ServerConfig{
			OzonCfg: &config.OzonConfig{
				BaseURL:             "baseurl",
				SearchBarSelector:   "searchbarselector",
				ItemsSelector:       "itemsselector",
				LinkSelector:        "linkselector",
				PriceSelector:       "priceselector",
				ProductNameSelector: "productnameselector",
				RatingSelector:      "ratingselector",
				ReviewsSelector:     "reviewsselector",
			},
		},
	}
	oz := parsers.NewOzonParser(cfg, loggerMock, browserRepoMock)

	t.Run("success", func(t *testing.T) {
		p := domain.Product{
			Name:         "product",
			Link:         cfg.Server.OzonCfg.BaseURL+"link",
			Price:        100.0,
			Rating:       5.0,
			ReviewsCount: 253,
		}

		prods := make([]domain.Product, 0, 1)
		prods = append(prods, p)

		prodName := "macbook pro 16gb 512gb"
		priceFrom := 50000.0
		priceTo := 250000.0

		linkPtr := "link"

		browserRepoMock.On("NewPage", mock.Anything).Return(pageMock, nil).Once()
		pageMock.On("Close").Return(nil).Once()
		pageMock.On("NavigateWithReferer", mock.Anything, cfg.Server.OzonCfg.BaseURL).Return(nil).Once()
		pageMock.On("WaitDOMStable", mock.Anything).Return(nil).Twice()

		pageMock.On("Element", mock.Anything, cfg.Server.OzonCfg.SearchBarSelector).Return(searchBarMock, nil).Once()
		pageMock.On("MoveCursorToElement", mock.Anything, cfg.Server.OzonCfg.SearchBarSelector).Return(nil).Once()
		searchBarMock.On("Click", mock.Anything).Return(nil).Once()
		searchBarMock.On("Input", mock.Anything, prodName).Return(nil).Once()
		pageMock.On("KeyboardType", mock.Anything, mock.Anything).Return(nil).Once()

		pageMock.On("Elements", mock.Anything, cfg.Server.OzonCfg.ItemsSelector).Return([]repository.Element{itemMock}, nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.LinkSelector).Return(linkElMock, nil).Once()
		linkElMock.On("Attribute", mock.Anything, "href").Return(&linkPtr, nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.PriceSelector).Return(priceElMock, nil).Once()
		priceElMock.On("Text", mock.Anything).Return("100 ₽", nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.ProductNameSelector).Return(prodNameElMock, nil).Once()
		prodNameElMock.On("Text", mock.Anything).Return("product", nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.RatingSelector).Return(ratingElMock, nil).Once()
		ratingElMock.On("Text", mock.Anything).Return("5.0", nil).Once()

		itemMock.On("ElementX", mock.Anything, cfg.Server.OzonCfg.ReviewsSelector).Return(reviewsElMock, nil).Once()
		reviewsElMock.On("Text", mock.Anything).Return("253", nil).Once()

		res, err := oz.GetAllProducts(context.Background(), prodName, priceFrom, priceTo)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.ElementsMatch(t, prods, res)

		browserRepoMock.AssertExpectations(t)
		pageMock.AssertExpectations(t)
		searchBarMock.AssertExpectations(t)
		itemMock.AssertExpectations(t)
		linkElMock.AssertExpectations(t)
		priceElMock.AssertExpectations(t)
		ratingElMock.AssertExpectations(t)
		reviewsElMock.AssertExpectations(t)
	})

	t.Run("zero rating/price/reviews", func(t *testing.T) {
		p := domain.Product{
			Name:         "product",
			Link:         cfg.Server.OzonCfg.BaseURL+"link",
			Price:        0.0,
			Rating:       0.0,
			ReviewsCount: 0,
		}

		prods := make([]domain.Product, 0, 1)
		prods = append(prods, p)

		prodName := "macbook pro 16gb 512gb"
		priceFrom := 50000.0
		priceTo := 250000.0

		linkPtr := "link"

		browserRepoMock.On("NewPage", mock.Anything).Return(pageMock, nil).Once()
		pageMock.On("Close").Return(nil).Once()
		pageMock.On("NavigateWithReferer", mock.Anything, cfg.Server.OzonCfg.BaseURL).Return(nil).Once()
		pageMock.On("WaitDOMStable", mock.Anything).Return(nil).Twice()

		pageMock.On("Element", mock.Anything, cfg.Server.OzonCfg.SearchBarSelector).Return(searchBarMock, nil).Once()
		pageMock.On("MoveCursorToElement", mock.Anything, cfg.Server.OzonCfg.SearchBarSelector).Return(nil).Once()
		searchBarMock.On("Click", mock.Anything).Return(nil).Once()
		searchBarMock.On("Input", mock.Anything, prodName).Return(nil).Once()
		pageMock.On("KeyboardType", mock.Anything, mock.Anything).Return(nil).Once()

		pageMock.On("Elements", mock.Anything, cfg.Server.OzonCfg.ItemsSelector).Return([]repository.Element{itemMock}, nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.LinkSelector).Return(linkElMock, nil).Once()
		linkElMock.On("Attribute", mock.Anything, "href").Return(&linkPtr, nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.PriceSelector).Return(priceElMock, nil).Once()
		loggerMock.On("Error", "parser string to float64 price", mock.Anything).Once()
		priceElMock.On("Text", mock.Anything).Return("invalid", nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.ProductNameSelector).Return(prodNameElMock, nil).Once()
		prodNameElMock.On("Text", mock.Anything).Return("product", nil).Once()

		itemMock.On("Element", mock.Anything, cfg.Server.OzonCfg.RatingSelector).Return(ratingElMock, nil).Once()
		loggerMock.On("Error", "parser string to float64 rating", mock.Anything).Once()
		ratingElMock.On("Text", mock.Anything).Return("invalid", nil).Once()

		itemMock.On("ElementX", mock.Anything, cfg.Server.OzonCfg.ReviewsSelector).Return(reviewsElMock, nil).Once()
		loggerMock.On("Error", "parser string to integer reviews", mock.Anything).Once()
		reviewsElMock.On("Text", mock.Anything).Return("invalid", nil).Once()

		res, err := oz.GetAllProducts(context.Background(), prodName, priceFrom, priceTo)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.ElementsMatch(t, prods, res)

		browserRepoMock.AssertExpectations(t)
		pageMock.AssertExpectations(t)
		searchBarMock.AssertExpectations(t)
		itemMock.AssertExpectations(t)
		linkElMock.AssertExpectations(t)
		priceElMock.AssertExpectations(t)
		ratingElMock.AssertExpectations(t)
		reviewsElMock.AssertExpectations(t)
	})

}
