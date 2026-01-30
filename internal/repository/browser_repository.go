package repository

import (
	"context"

	"github.com/go-rod/rod"
)

type BrowserRepository interface {
	Connect(ctx context.Context) (*rod.Browser, error)
	Ping(ctx context.Context) error

	NewPage(ctx context.Context) (Page, error)
}