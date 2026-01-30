package chromium_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/browser/chromium"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
)

func TestChromium_NewBrowser(t *testing.T) {
	cfg := fakeConfig(t)

	ch := chromium.NewChromiumRepository(cfg)
	assert.NotNil(t, ch)
	browser := chromium.NewBrowser(ch)
	assert.NotNil(t, browser)
}

func TestChromium_Chromium(t *testing.T) {
	cfg := fakeConfig(t)

	ch := chromium.NewChromiumRepository(cfg)
	assert.NotNil(t, ch)
	browser := chromium.NewBrowser(ch)
	assert.NotNil(t, browser)
	rep := browser.Chromium()
	assert.NotNil(t, rep)
}

func fakeConfig(t *testing.T) *config.Config {
	t.Helper()

	return &config.Config{
		Browser: config.BrowserConfig{
			WsURL: "ws://1.2.3.4:3000/chromium",
		},
		Options: config.OptionsConfig{
			LoggerTimeFormat: "02-01-2006 15:04:05",
		},
		Server: config.ServerConfig{
			Env:      "local",
			HTTPAddr: "localhost:8080",
		},
	}
}
