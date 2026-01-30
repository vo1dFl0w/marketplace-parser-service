package http_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	ht "github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http/httpgen"
)

func TestErrors_Error(t *testing.T) {
	httpErr := &ht.HTTPError{Message: "msg", Status: http.StatusOK}

	res := httpErr.Error()
	assert.NotNil(t, res)
	assert.Equal(t, httpErr.Message, res)
}

func TestErrors_ToSearchProductErrResp(t *testing.T) {
	testCases := []struct {
		name    string
		httpErr *ht.HTTPError
	}{
		{
			name:    "Bad Request",
			httpErr: &ht.HTTPError{Message: "msg", Status: http.StatusBadRequest},
		},
		{
			name:    "Client Closed Request",
			httpErr: &ht.HTTPError{Message: "msg", Status: ht.StatusClientClosedRequest},
		},
		{
			name:    "Internal Server Error",
			httpErr: &ht.HTTPError{Message: "msg", Status: http.StatusInternalServerError},
		},
		{
			name:    "Gateway Timeout",
			httpErr: &ht.HTTPError{Message: "msg", Status: http.StatusGatewayTimeout},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.httpErr.ToSearchProductErrResp()
			assert.NotNil(t, res)

			switch tc.httpErr.Status {
			case http.StatusBadRequest:
				_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetBadRequest)
				assert.True(t, ok)
			case ht.StatusClientClosedRequest:
				_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetCode499)
				assert.True(t, ok)
			case http.StatusInternalServerError:
				_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetInternalServerError)
				assert.True(t, ok)
			case http.StatusGatewayTimeout:
				_, ok := res.(*httpgen.APIV1MarketplaceParserServiceProductsSearchGetGatewayTimeout)
				assert.True(t, ok)
			}
		})
	}
}
