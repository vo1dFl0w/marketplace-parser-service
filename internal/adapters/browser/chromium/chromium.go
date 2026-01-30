package chromium

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/config"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
)

const (
	userAgentWindows = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	platformWindows  = "Win32"

	// userAgentLinux = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	// platformLinux  = "Linux x86_64"

	// userAgentMacOS = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	// platoformMacOS = "MacIntel"
)

type ChromiumRepository struct {
	cfg *Config
}

// Create a new chromium repository struct.
func NewChromiumRepository(cfg *config.Config) *ChromiumRepository {
	return &ChromiumRepository{cfg: NewChromiumConfig(cfg)}
}

// Pings browser
func (r *ChromiumRepository) Ping(ctx context.Context) error {
	ctxBrowser, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	browser, err := r.Connect(ctxBrowser)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return repository.ErrGatewayTimeout
		} else if errors.Is(err, context.Canceled) {
			return repository.ErrClientClosedRequest
		} else {
			return err
		}
	}
	defer browser.Close()

	return nil
}

// Connects to browser.
func (r *ChromiumRepository) Connect(ctx context.Context) (*rod.Browser, error) {

	// ====== docker container (headless=true) ======
	l, err := launcher.NewManaged(r.cfg.wsURL)
	if err != nil {
		return nil, fmt.Errorf("new manager: %w", err)
	}

	// Launcher configuration to bypass anti-fraud systems (Ozon/WB) and optimize for Docker:
	// - user-agent: emulate a real user.
	// - headless=new: use a modern engine with full support for graphics APIs.
	// - AutomationControlled: hide the navigator.webdriver flag.
	// - disable-dev-shm-usage: use RAM instead of /dev/shm to avoid crashes in Docker.
	l.Set("user-agent", userAgentWindows).
		Set("headless", "new").
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-infobars").
		Set("disable-dev-shm-usage").
		Set("lang", "ru-RU").
		Set("disable-features", "IsolateOrigins,site-per-process")

	c, err := l.Client()
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	browser := rod.New().Timeout(time.Second * 30).Client(c).Context(ctx)
	if err := browser.Connect(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, repository.ErrGatewayTimeout
		} else if errors.Is(err, context.Canceled) {
			return nil, repository.ErrClientClosedRequest
		} else {
			return nil, fmt.Errorf("connect chromium: %w", err)
		}
	}

	// ===== debug browser (headless=false) ======
	/*
		u := launcher.
			New().Headless(false).
			Set("disable-blink-features", "AutomationControlled").
			Set("disable-infobars").
			Set("disable-dev-shm-usage").
			Set("lang", "ru-RU").
			Set("disable-features", "IsolateOrigins,site-per-process").
			MustLaunch()

		// Trace enables/disables the visual tracing of the input actions on the page
		browser := rod.New().ControlURL(u).Trace(false).MustConnect()

	*/

	return browser, nil
}

// Creates new page with flags and user-agent
func (r *ChromiumRepository) NewPage(ctx context.Context) (repository.Page, error) {
	browser, err := r.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect browser: %w", err)
	}

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, repository.ErrGatewayTimeout
		} else if errors.Is(err, context.Canceled) {
			return nil, repository.ErrClientClosedRequest
		} else {
			return nil, fmt.Errorf("page: %w", err)
		}
	}

	if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      userAgentWindows,
		AcceptLanguage: r.cfg.acceptLanguage,
		Platform:       platformWindows,
	}); err != nil {
		return nil, fmt.Errorf("set user agent: %w", err)
	}

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             1920,
		Height:            1080,
		DeviceScaleFactor: 1,
		Mobile:            false,
	}); err != nil {
		return nil, fmt.Errorf("set view port: %w", err)
	}

	// Debug the navigator.webdriver value, if the webdriver=true, then the browser launch arguments were not set.
	// webdriver, _ := page.Evaluate(&rod.EvalOptions{JS: `() => navigator.webdriver`})
	// fmt.Printf("navigator.webdriver=%v\n", webdriver)

	// Debug the navigator.userAgent value.
	// ua, err := page.Evaluate(&rod.EvalOptions{JS: `() => navigator.userAgent`})
	// fmt.Printf("navigator.userAgent=%v\n", ua)

	return &rodPage{
		page:    page,
		browser: browser,
		cfg:     r.cfg,
	}, nil
}
