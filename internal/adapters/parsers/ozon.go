package parsers

import (
	"context"
	"fmt"
	"math/rand"

	"time"

	"github.com/go-rod/rod/lib/input"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/logger"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/utils"
)

type OzonParser interface {
	GetAllProducts(ctx context.Context, name string, priceFrom float64, priceTo float64) ([]domain.Product, error)
}

type ozonParser struct {
	cfg     *OzonConfig
	logger  logger.Logger
	browser repository.BrowserRepository
}

func NewOzonParser(cfg *config.Config, logger logger.Logger, browser repository.BrowserRepository) *ozonParser {
	return &ozonParser{cfg: NewOzonConfig(cfg), logger: logger, browser: browser}
}

func (op *ozonParser) GetAllProducts(ctx context.Context, name string, priceFrom float64, priceTo float64) ([]domain.Product, error) {
	page, err := op.browser.NewPage(ctx)
	if err != nil {
		return nil, utils.WrapError("page", err, ctx)
	}
	defer page.Close()

	if err := page.NavigateWithReferer(ctx, op.cfg.BaseURL); err != nil {
		return nil, utils.WrapError("navigate page with referer", err, ctx)
	}

	// Wait for the DOM to load to find the closing button.
	if err := page.WaitDOMStable(ctx); err != nil {
		return nil, utils.WrapError("wait dom stable", err, ctx)
	}

	searchBar, err := page.Element(ctx, op.cfg.SearchBarSelector)
	if err != nil {
		return nil, utils.WrapError("element search bar", err, ctx)
	}

	if err := page.MoveCursorToElement(ctx, op.cfg.SearchBarSelector); err != nil {
		return nil, utils.WrapError("move cursor to element search bar", err, ctx)
	}

	if err := searchBar.Click(ctx); err != nil {
		return nil, utils.WrapError("click search bar", err, ctx)
	}
	time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)

	if err := searchBar.Input(ctx, name); err != nil {
		return nil, utils.WrapError("input search bar", err, ctx)
	}
	time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)

	if err := page.KeyboardType(ctx, input.Enter); err != nil {
		return nil, utils.WrapError("type enter", err, ctx)
	}

	// Wait for the DOM to load to find the closing button.
	if err := page.WaitDOMStable(ctx); err != nil {
		return nil, utils.WrapError("wait dom stable", err, ctx)
	}

	items, err := page.Elements(ctx, op.cfg.ItemsSelector)
	if err != nil {
		return nil, utils.WrapError("elemets", err, ctx)
	}

	res := make([]domain.Product, 0, len(items))

	for _, itm := range items {
		if len(res) == 10 {
			break
		}

		var href string
		itmLink, _ := itm.Element(ctx, op.cfg.LinkSelector)
		if itmLink != nil {
			hrefRaw, err := itmLink.Attribute(ctx, "href")
			if err != nil {
				return nil, utils.WrapError("attribute link", err, ctx)
			}
			href = fmt.Sprintf("%s%s", op.cfg.BaseURL, *hrefRaw)
		}

		var price float64
		itmPrice, _ := itm.Element(ctx, op.cfg.PriceSelector)
		if itmPrice != nil {
			priceStr, _ := itmPrice.Text(ctx)
			price, err = ParseStringToFloat64(priceStr)
			if err != nil {
				op.logger.Error("parser string to float64 price", err)
				price = 0.0
				// return nil, utils.WrapError("parse string to float64 price", err, ctx)
			}
		}

		var name string
		itmName, _ := itm.Element(ctx, op.cfg.ProductNameSelector)
		if itmName != nil {
			name, _ = itmName.Text(ctx)
		}

		if href == "" || name == "" {
			continue
		}

		var rating float64
		itmRating, _ := itm.Element(ctx, op.cfg.RatingSelector)
		if itmRating != nil {
			ratingStr, _ := itmRating.Text(ctx)
			rating, err = ParseStringToFloat64(ratingStr)
			if err != nil {
				op.logger.Error("parser string to float64 rating", err)
				rating = 0.0
				// return nil, WrapError("parse string to float64 rating", err, ctx)
			}
		}

		var reviews int
		itmReviews, _ := itm.ElementX(ctx, op.cfg.ReviewsSelector)
		if itmReviews != nil {
			reviewsStr, _ := itmReviews.Text(ctx)
			reviews, err = ParseStringToInteger(reviewsStr)
			if err != nil {
				op.logger.Error("parser string to integer reviews", err)
				reviews = 0
				// return nil, WrapError("parse string to integer reviews", err, ctx)
			}
		}

		res = append(res, domain.Product{
			Name:         name,
			Link:         href,
			Price:        price,
			Rating:       rating,
			ReviewsCount: reviews,
		})
	}

	return res, nil
}
