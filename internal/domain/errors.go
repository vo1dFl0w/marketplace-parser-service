package domain

import "errors"

var (
	ErrEmptyProductName      = errors.New("empty product name")
	ErrGatewayTimeout        = errors.New("gateway timeout")
	ErrClientClosedRequest   = errors.New("client closed request")
	ErrPriceFromBelowZero    = errors.New("price from below zero")
	ErrPriceFromAbovePriceTo = errors.New("price from above price to")
	ErrPriceToBelowZero      = errors.New("price to below zero")
)
