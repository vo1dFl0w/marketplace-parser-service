package parsers

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-rod/rod/lib/input"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/logger"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/utils"
)

type WildberriesParser interface {
	GetAllProducts(ctx context.Context, name string, priceFrom float64, priceTo float64) ([]domain.Product, error)
}

type wildberriesParser struct {
	cfg     *WildberriesConfig
	logger  logger.Logger
	browser repository.BrowserRepository
}

// NewWildberriesParser сreate a new empty object that implements the WildberriesParser interface.
func NewWildberriesParser(cfg *config.Config, logger logger.Logger, browser repository.BrowserRepository) *wildberriesParser {
	return &wildberriesParser{cfg: NewWildberriesConfig(cfg), logger: logger, browser: browser}
}

// GetAllProducts parses and gets a list of products from the site.
func (wp *wildberriesParser) GetAllProducts(ctx context.Context, name string, priceFrom float64, priceTo float64) ([]domain.Product, error) {
	// Open stealth page
	page, err := wp.browser.NewPage(ctx)
	if err != nil {
		return nil, utils.WrapError("page", err, ctx)
	}
	defer page.Close()

	// Navigate to wb by base url from config
	if err := page.NavigateWithReferer(ctx, wp.cfg.BaseURL); err != nil {
		return nil, utils.WrapError("navigate page with referer", err, ctx)
	}
	// Wait for the DOM to load to find the closing button.
	if err := page.WaitDOMStable(ctx); err != nil {
		return nil, utils.WrapError("wait dom stable", err, ctx)
	}

	// Close the pop-up window if there is one
	if err := page.ClosePopUpWindow(ctx, wp.cfg.CloseButtonSelector); err != nil {
		return nil, utils.WrapError("close pop up window", err, ctx)
	}

	// Find wb serach bar
	searchBar, err := page.Element(ctx, wp.cfg.SearchBarSelector)
	if err != nil {
		return nil, utils.WrapError("element search bar", err, ctx)
	}
	// Simulates moving the cursor across the page to an object
	if err := page.MoveCursorToElement(ctx, wp.cfg.SearchBarSelector); err != nil {
		return nil, utils.WrapError("move cursor to element search bar", err, ctx)
	}
	if err := searchBar.Click(ctx); err != nil {
		return nil, utils.WrapError("click search bar", err, ctx)
	}
	time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)

	/*for _, r := range name {
		// Simulates typing characters in a search bar
		if err := searchBar.Input(string(r)); err != nil {
			return nil, WrapError("input search bar", err, ctx)
		}
		time.Sleep(100 * time.Millisecond)
	}
	*/

	if err := searchBar.Input(ctx, name); err != nil {
		return nil, utils.WrapError("input search bar", err, ctx)
	}
	time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
	if err := page.KeyboardType(ctx, input.Enter); err != nil {
		return nil, utils.WrapError("type enter", err, ctx)
	}
	// Wait for DOM to load for further parsing of product cards
	if err := page.WaitDOMStable(ctx); err != nil {
		return nil, utils.WrapError("wait dom stable", err, ctx)
	}

	// Find and parse product-cards and parse
	items, err := page.Elements(ctx, wp.cfg.ItemsSelector)
	if err != nil {
		return nil, utils.WrapError("find elements", err, ctx)
	}

	res := make([]domain.Product, 0, len(items))
	for _, itm := range items {
		if len(res) == 10 {
			break
		}
		// Find product name and link
		var link string
		var name string
		itmLink, _ := itm.Element(ctx, wp.cfg.LinkSelector)
		if itmLink != nil {
			href, err := itmLink.Attribute(ctx, "href")
			if err != nil {
				return nil, utils.WrapError("attribute link", err, ctx)
			}
			if href != nil {
				link = *href
			}
			label, err := itmLink.Attribute(ctx, "aria-label")
			if err != nil {
				return nil, utils.WrapError("attribute label", err, ctx)
			}
			if label != nil {
				name = *label
			}
		}
		// Find product price
		var price float64
		itmPrice, _ := itm.Element(ctx, wp.cfg.PriceSelector)
		if itmPrice != nil {
			priceStr, err := itmPrice.Text(ctx)
			if err != nil {
				return nil, utils.WrapError("text price", err, ctx)
			}
			price, err = ParseStringToFloat64(priceStr)
			if err != nil {
				// Log if an error occurs while parsing string to float64 and set price = 0
				wp.logger.Error("parser string to float64 price", err)
				price = 0.0
				// return nil, WrapError("parse string to float64 price", err, ctx)
			}
		}

		if link == "" || name == "" {
			continue
		}
		// Find product rating
		var rating float64
		itmRating, _ := itm.Element(ctx, wp.cfg.RatingSelector)
		if itmRating != nil {
			ratingStr, err := itmRating.Text(ctx)
			if err != nil {
				return nil, utils.WrapError("text rating", err, ctx)
			}
			rating, err = ParseStringToFloat64(ratingStr)
			if err != nil {
				// Log if an error occurs while parsing string to float64 and set rating = 0
				wp.logger.Error("parser string to float64 rating", err)
				rating = 0.0
				// return nil, WrapError("parse string to float64 rating", err, ctx)
			}
		}
		// Find product reviews
		var reviews int
		itmReviews, _ := itm.Element(ctx, wp.cfg.ReviewsSelector)
		if itmReviews != nil {
			reviewsStr, err := itmReviews.Text(ctx)
			if err != nil {
				return nil, utils.WrapError("text reviews", err, ctx)
			}
			reviews, err = ParseStringToInteger(reviewsStr)
			if err != nil {
				// Log if an error occurs while parsing string to integer and set reviews = 0
				wp.logger.Error("parser string to integer reviews", err)
				reviews = 0
				// return nil, WrapError("parse string to integer reviews", err, ctx)
			}
		}

		res = append(res, domain.Product{
			Name:         name,
			Link:         link,
			Price:        price,
			Rating:       rating,
			ReviewsCount: reviews,
		})
	}

	return res, nil
}
