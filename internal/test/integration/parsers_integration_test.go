package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/browser/chromium"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/parsers"
	integr "github.com/vo1dFl0w/marketplace-parser-service/internal/test/integration"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/logger"
)

func TestIntegration_WildberriesParser(t *testing.T) {
	timeout := time.Second * 30

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	prodName := "macbook pro 16gb 512gb"
	priceFrom := 50000.0
	priceTo := 250000.0

	loggerCfg := logger.NewLoggerConfig(integr.Cfg.Server.Env, integr.Cfg.Options.LoggerTimeFormat)
	logger := logger.LoadLogger(loggerCfg)

	chromiumRepo := chromium.NewChromiumRepository(integr.Cfg)
	browserRepo := chromium.NewBrowser(chromiumRepo)

	wb := parsers.NewWildberriesParser(integr.Cfg, logger, browserRepo.Chromium())

	res, err := wb.GetAllProducts(ctx, prodName, priceFrom, priceTo)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestIntegration_OzonParser(t *testing.T) {
	timeout := time.Second * 30

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	prodName := "macbook pro 16gb 512gb"
	priceFrom := 50000.0
	priceTo := 250000.0

	loggerCfg := logger.NewLoggerConfig(integr.Cfg.Server.Env, integr.Cfg.Options.LoggerTimeFormat)
	logger := logger.LoadLogger(loggerCfg)

	chromiumRepo := chromium.NewChromiumRepository(integr.Cfg)
	browserRepo := chromium.NewBrowser(chromiumRepo)

	oz := parsers.NewOzonParser(integr.Cfg, logger, browserRepo.Chromium())

	res, err := oz.GetAllProducts(ctx, prodName, priceFrom, priceTo)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	t.Log(res)
}
