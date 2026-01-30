package chromium

import (
	"time"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
)

type Config struct {
	wsURL          string
	referer        string
	acceptLanguage string
	// domStableDuration in milliseconds
	domStableDuration time.Duration
	domStableDiff     float64
}

func NewChromiumConfig(cfg *config.Config) *Config {
	return &Config{
		wsURL:             cfg.Browser.WsURL,
		referer:           cfg.Browser.Referer,
		acceptLanguage:    cfg.Browser.AcceptLanguage,
		domStableDuration: cfg.Browser.DomStableDuration,
		domStableDiff:     cfg.Browser.DomStableDiff,
	}
}
