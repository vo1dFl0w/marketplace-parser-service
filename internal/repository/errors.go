package repository

import "errors"

var (
	ErrGatewayTimeout      = errors.New("gateway timeout")
	ErrClientClosedRequest = errors.New("client closed request")
)
