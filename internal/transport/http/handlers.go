package http

import (
	"context"
	"net/http"
	"time"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/usecase"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/logger"
)

type Handler struct {
	logger         logger.Logger
	router         *http.ServeMux
	parserSrv      usecase.ParserService
	requestTimeout time.Duration
}

func NewHandler(logger logger.Logger, parserSrv usecase.ParserService, requestTimeout time.Duration) *Handler {
	return &Handler{logger: logger, router: http.NewServeMux(), parserSrv: parserSrv, requestTimeout: requestTimeout}
}

func (h *Handler) APIV1MarketplaceParserServiceProductsSearchGet(ctx context.Context, params httpgen.APIV1MarketplaceParserServiceProductsSearchGetParams) (httpgen.APIV1MarketplaceParserServiceProductsSearchGetRes, error) {
	prods, err := h.parserSrv.GetProductsList(ctx, params.Name, params.PriceFrom.Value, params.PriceTo.Value)
	if err != nil {
		httpErr := MapError(err)
		h.LogHTTPError(ctx, err, httpErr)
		return httpErr.ToSearchProductErrResp(), nil
	}
	res := make(httpgen.SearchProductsResponse, 0, len(prods))

	for _, p := range prods {
		res = append(res, httpgen.Product{
			Name:         p.Name,
			Link:         p.Link,
			Price:        p.Price,
			Rating:       p.Rating,
			ReviewsCount: p.ReviewsCount,
		})
	}

	return &res, nil
}

func (h *Handler) LogHTTPError(ctx context.Context, err error, httpErr *HTTPError) {
	attrs := []any{
		"error", err,
		"status", httpErr.Status,
		"message", httpErr.Message,
	}

	switch {
	case httpErr.Status >= 500:
		switch httpErr.Status {
		case http.StatusGatewayTimeout:
			h.logger.Error("http_request_failed", append(attrs, "reason", "dependency_timeout")...)
		default:
			h.logger.Error("http_request_failed", append(attrs, "reason", "internal_server_error")...)
		}
	case httpErr.Status >= 400:
		h.logger.Warn("http_request_failed", append(attrs, "reason", "client_error")...)
	}
}
