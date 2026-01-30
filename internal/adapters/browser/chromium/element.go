package chromium

import (
	"context"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
)

type rodElement struct {
	element *rod.Element
}

func (e *rodElement) Text(ctx context.Context) (string, error) {
	return e.element.Context(ctx).Text()
}

func (e *rodElement) Attribute(ctx context.Context, name string) (*string, error) {
	return e.element.Context(ctx).Attribute(name)
}

func (e *rodElement) Click(ctx context.Context) error {
	return e.element.Context(ctx).Click(proto.InputMouseButtonLeft, 1)
}

func (e *rodElement) Input(ctx context.Context, text string) error {
	return e.element.Context(ctx).Input(text)
}

func (e *rodElement) Element(ctx context.Context, selector string) (repository.Element, error) {
	el, err := e.element.Context(ctx).Element(selector)
	if err != nil {
		return nil, err
	}

	return &rodElement{
		element: el,
	}, nil
}

func (e *rodElement) ElementX(ctx context.Context, selector string) (repository.Element, error) {
	elem, err := e.element.Context(ctx).ElementX(selector)
	if err != nil {
		return nil, err
	}

	return &rodElement{element: elem}, nil
}
