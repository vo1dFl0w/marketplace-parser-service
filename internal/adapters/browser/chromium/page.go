package chromium

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
)

type rodPage struct {
	page    *rod.Page
	browser *rod.Browser
	cfg     *Config
}

// NavigatePageWithReferer navigates current page to the given baseURL.
func (p *rodPage) NavigateWithReferer(ctx context.Context, baseURL string) error {
	_, err := proto.PageNavigate{
		URL:      baseURL,
		Referrer: p.cfg.referer,
	}.Call(p.page)
	if err != nil {
		return err
	}

	return nil
}

// WaitDOMStable waits until the change of the DOM tree is less or equal than domStableDiff percent for domStableDuration.
func (p *rodPage) WaitDOMStable(ctx context.Context) error {
	if err := p.page.Context(ctx).WaitDOMStable(p.cfg.domStableDuration, p.cfg.domStableDiff); err != nil {
		return err
	}

	return nil
}

// ClosePopUpWindow tries to find a close button on pop-up window and click on it
func (p *rodPage) ClosePopUpWindow(ctx context.Context, selector string) error {
	b, el, err := p.page.Has(selector)
	if err != nil {
		return err
	}
	if b {
		// Simulates moving the cursor across the page to an object
		if err := p.MoveCursorToElement(ctx, selector); err != nil {
			return err
		}
		if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
			return err
		}
		time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
	}

	return nil
}

// Element retries until an element in the page that matches the CSS selector, then returns the matched element.
func (p *rodPage) Element(ctx context.Context, selector string) (repository.Element, error) {
	elem, err := p.page.Context(ctx).Element(selector)
	if err != nil {
		return nil, err
	}

	return &rodElement{
		element: elem,
	}, nil
}

// Elements returns all elements that match the css selector.
func (p *rodPage) Elements(ctx context.Context, selector string) ([]repository.Element, error) {
	elems, err := p.page.Context(ctx).Elements(selector)
	if err != nil {
		return nil, err
	}

	res := make([]repository.Element, 0, len(elems))

	for _, e := range elems {
		res = append(res, &rodElement{element: e})
	}

	return res, nil
}

// KeyBoardType simulates pressing a key.
func (p *rodPage) KeyboardType(ctx context.Context, key input.Key) error {
	if err := p.page.Keyboard.Type(key); err != nil {
		return err
	}

	return nil
}

// MoveCursorToElement simulates moving the mouse cursor on the browser page to a specified element
func (p *rodPage) MoveCursorToElement(ctx context.Context, elemName string) error {
	elem, err := p.page.Context(ctx).Element(elemName)
	if err != nil {
		return err
	}

	shape, err := elem.Shape()
	if err != nil {
		return err
	}
	box := shape.Box()

	// Find the center of the box
	centerX := box.X + box.Width/2
	centerY := box.Y + box.Height/2
	// Simulate move mouse to the center of the box
	if err := p.page.Mouse.MoveLinear(proto.Point{X: centerX, Y: centerY}, 15); err != nil {
		return err
	}

	return nil
}

// Close closes active page and browser
func (p *rodPage) Close() error {
	if p.page != nil {
		if err := p.page.Close(); err != nil {
			return err
		}
	}
	if p.browser != nil {
		if err := p.browser.Close(); err != nil {
			return err
		}
	}

	return nil
}
