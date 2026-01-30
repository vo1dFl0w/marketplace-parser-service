package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/browser/chromium"
	integr "github.com/vo1dFl0w/marketplace-parser-service/internal/test/integration"
)

func TestIntegration_ChromiumConnect(t *testing.T) {
	chromiumRepo := chromium.NewChromiumRepository(integr.Cfg)
	b := chromium.NewBrowser(chromiumRepo)

	browser, err := b.Chromium().Connect(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, browser)

	err = browser.Close()
	assert.NoError(t, err)
}

func TestIntegration_ChromiumPing(t *testing.T) {
	chromiumRepo := chromium.NewChromiumRepository(integr.Cfg)
	browser := chromium.NewBrowser(chromiumRepo)

	err := browser.Chromium().Ping(context.Background())
	assert.NoError(t, err)
}
