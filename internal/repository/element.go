package repository

import "context"

type Element interface {
	Text(ctx context.Context) (string, error)
	Attribute(ctx context.Context, name string) (*string, error)
	Click(ctx context.Context) error
	Input(ctx context.Context, text string) error
	Element(ctx context.Context, selector string) (Element, error)
	ElementX(ctx context.Context, selector string) (Element, error)
}