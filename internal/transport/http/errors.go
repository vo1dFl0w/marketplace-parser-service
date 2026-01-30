package http

import (
	"errors"
	"net/http"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/domain"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http/httpgen"
)

const (
	StatusClientClosedRequest = 499
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrClientClosedRequest = errors.New("client closed request")
	ErrGatewayTimeout      = errors.New("gateway timeout")
	ErrInternalServerError = errors.New("internal server error")
)

type HTTPError struct {
	Message string
	Status  int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) ToSearchProductErrResp() httpgen.APIV1MarketplaceParserServiceProductsSearchGetRes {
	switch e.Status {
	case http.StatusBadRequest:
		return &httpgen.APIV1MarketplaceParserServiceProductsSearchGetBadRequest{Message: e.Message, Status: e.Status}
	case StatusClientClosedRequest:
		return &httpgen.APIV1MarketplaceParserServiceProductsSearchGetCode499{Message: e.Message, Status: e.Status}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1MarketplaceParserServiceProductsSearchGetGatewayTimeout{Message: e.Message, Status: e.Status}
	default:
		return &httpgen.APIV1MarketplaceParserServiceProductsSearchGetInternalServerError{Message: e.Message, Status: e.Status}
	}
}

func MapError(err error) *HTTPError {
	switch {
	case errors.Is(err, domain.ErrEmptyProductName):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrClientClosedRequest):
		return &HTTPError{Message: ErrClientClosedRequest.Error(), Status: StatusClientClosedRequest}
	case errors.Is(err, domain.ErrPriceFromBelowZero):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrPriceFromAbovePriceTo):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrPriceToBelowZero):
		return &HTTPError{Message: ErrBadRequest.Error(), Status: http.StatusBadRequest}
	case errors.Is(err, domain.ErrGatewayTimeout):
		return &HTTPError{Message: ErrGatewayTimeout.Error(), Status: http.StatusGatewayTimeout}
	default:
		return &HTTPError{Message: ErrInternalServerError.Error(), Status: http.StatusInternalServerError}
	}
}
