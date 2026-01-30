package utils_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
	"github.com/vo1dFl0w/marketplace-parser-service/pkg/utils"
)

func TestParses_MapContextOnly(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		ctx    context.Context
		expErr error
	}{
		{
			name:   "any error",
			err:    errors.New("any error"),
			ctx:    context.Background(),
			expErr: nil,
		},
		{
			name:   "deadline exceeded error",
			err:    context.DeadlineExceeded,
			ctx:    context.Background(),
			expErr: repository.ErrGatewayTimeout,
		},
		{
			name:   "client closed request",
			err:    context.Canceled,
			ctx:    context.Background(),
			expErr: repository.ErrClientClosedRequest,
		},
		{
			name: "context deadline exceeded error",
			err:  context.DeadlineExceeded,
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				cancel()
				return ctx
			}(),
			expErr: repository.ErrGatewayTimeout,
		},
		{
			name: "context canceled",
			err:  context.Canceled,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expErr: repository.ErrClientClosedRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expErr != nil {
				err := utils.MapContextOnly(tc.err, tc.ctx)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expErr)
			} else {
				err := utils.MapContextOnly(tc.err, tc.ctx)
				assert.NoError(t, err)
			}
		})
	}
}

func TestParses_WrapError(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		ctx    context.Context
		expErr error
	}{
		{
			name:   "any error",
			err:    errors.New("any error"),
			ctx:    context.Background(),
			expErr: errors.New("any error"),
		},
		{
			name:   "deadline exceeded error",
			err:    context.DeadlineExceeded,
			ctx:    context.Background(),
			expErr: repository.ErrGatewayTimeout,
		},
		{
			name:   "client closed request",
			err:    context.Canceled,
			ctx:    context.Background(),
			expErr: repository.ErrClientClosedRequest,
		},
		{
			name: "context deadline exceeded error",
			err:  context.DeadlineExceeded,
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				cancel()
				return ctx
			}(),
			expErr: repository.ErrGatewayTimeout,
		},
		{
			name: "context canceled",
			err:  context.Canceled,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expErr: repository.ErrClientClosedRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.WrapError("error", tc.err, tc.ctx)
			assert.Error(t, err)
			assert.ErrorContains(t, err, tc.expErr.Error())
		})
	}
}
