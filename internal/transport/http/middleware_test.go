package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/test/mocks"
	ht "github.com/vo1dFl0w/marketplace-parser-service/internal/transport/http"
)

func TestMiddlewares_LoggerMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		status        int
		expLevel      string
		expStatusText string
	}{
		{
			name:          "200",
			status:        http.StatusOK,
			expLevel:      "info",
			expStatusText: http.StatusText(http.StatusOK),
		},
		{
			name:          "400",
			status:        http.StatusBadRequest,
			expLevel:      "warn",
			expStatusText: http.StatusText(http.StatusBadRequest),
		},
		{
			name:          "499",
			status:        ht.StatusClientClosedRequest,
			expLevel:      "warn",
			expStatusText: ht.ErrClientClosedRequest.Error(),
		},
		{
			name:          "500",
			status:        http.StatusInternalServerError,
			expLevel:      "error",
			expStatusText: http.StatusText(http.StatusInternalServerError),
		},
		{
			name:          "504",
			status:        http.StatusGatewayTimeout,
			expLevel:      "error",
			expStatusText: http.StatusText(http.StatusGatewayTimeout),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timeout := time.Second * 30
			loggerMock := &mocks.LoggerMock{}

			handler := ht.NewHandler(loggerMock, nil, timeout)

			req := httptest.NewRequest(http.MethodGet, "/testmiddleware", nil)
			req.RemoteAddr = "1.2.3.4:1234"
			req.Method = http.MethodGet

			loggerMock.On("With", []any{"remote_addr", req.RemoteAddr, "http-method", req.Method, "path", req.URL.Path}).Return(loggerMock).Once()
			loggerMock.On("Info", "started").Once()

			switch {
			case tc.status >= 500:
				loggerMock.On("Error", "failed", mock.Anything).Once()
			case tc.status >= 400:
				loggerMock.On("Warn", "failed", mock.Anything).Once()
			default:
				loggerMock.On("Info", "completed", mock.Anything).Once()
			}

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte("ok"))
			})

			middleware := handler.LoggerMiddleware(next)

			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			assert.Equal(t, tc.status, rec.Code)
			assert.Equal(t, "ok", rec.Body.String())

			loggerMock.AssertExpectations(t)
		})
	}
}
