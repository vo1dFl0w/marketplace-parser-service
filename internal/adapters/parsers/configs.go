package parsers

import "github.com/vo1dFl0w/marketplace-parser-service/internal/config"

type WildberriesConfig struct {
	BaseURL             string
	CloseButtonSelector string
	SearchBarSelector   string
	ItemsSelector       string
	LinkSelector        string
	PriceSelector       string
	RatingSelector      string
	ReviewsSelector     string
}

func NewWildberriesConfig(cfg *config.Config) *WildberriesConfig {
	return &WildberriesConfig{
		BaseURL:             cfg.Server.WbCfg.BaseURL,
		CloseButtonSelector: cfg.Server.WbCfg.CloseButtonSelector,
		SearchBarSelector:   cfg.Server.WbCfg.SearchBarSelector,
		ItemsSelector:       cfg.Server.WbCfg.ItemsSelector,
		LinkSelector:        cfg.Server.WbCfg.LinkSelector,
		PriceSelector:       cfg.Server.WbCfg.PriceSelector,
		RatingSelector:      cfg.Server.WbCfg.RatingSelector,
		ReviewsSelector:     cfg.Server.WbCfg.ReviewsSelector,
	}
}

type OzonConfig struct {
	BaseURL             string
	SearchBarSelector   string
	ItemsSelector       string
	LinkSelector        string
	ProductNameSelector string
	PriceSelector       string
	RatingSelector      string
	ReviewsSelector     string
}

func NewOzonConfig(cfg *config.Config) *OzonConfig {
	return &OzonConfig{
		BaseURL:             cfg.Server.OzonCfg.BaseURL,
		SearchBarSelector:   cfg.Server.OzonCfg.SearchBarSelector,
		ItemsSelector:       cfg.Server.OzonCfg.ItemsSelector,
		LinkSelector:        cfg.Server.OzonCfg.LinkSelector,
		PriceSelector:       cfg.Server.OzonCfg.PriceSelector,
		ProductNameSelector: cfg.Server.OzonCfg.ProductNameSelector,
		RatingSelector:      cfg.Server.OzonCfg.RatingSelector,
		ReviewsSelector:     cfg.Server.OzonCfg.ReviewsSelector,
	}
}
