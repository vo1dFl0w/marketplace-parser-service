package repository

import (
	"context"

	"github.com/go-rod/rod/lib/input"
)

type Page interface {
	NavigateWithReferer(ctx context.Context, url string) error
	WaitDOMStable(ctx context.Context) error
	ClosePopUpWindow(ctx context.Context, selector string) error
	MoveCursorToElement(ctx context.Context, elemName string) error
	Element(ctx context.Context, selector string) (Element, error)
	Elements(ctx context.Context, selector string) ([]Element, error)
	KeyboardType(ctx context.Context, key input.Key) error
	Close() error
}
