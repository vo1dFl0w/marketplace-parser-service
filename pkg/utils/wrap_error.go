package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/vo1dFl0w/marketplace-parser-service/internal/repository"
)

func MapContextOnly(err error, ctx context.Context) error {
    if errors.Is(err, context.DeadlineExceeded) {
        return repository.ErrGatewayTimeout
    }
    if errors.Is(err, context.Canceled) {
        return repository.ErrClientClosedRequest
    }

    if ctx != nil {
        if errCtx := ctx.Err(); errCtx != nil {
            if errors.Is(errCtx, context.DeadlineExceeded) {
                return repository.ErrGatewayTimeout
            }
            if errors.Is(errCtx, context.Canceled) {
                return repository.ErrClientClosedRequest
            }
        }
    }

    return nil
}

func WrapError(prefix string, err error, ctx context.Context) error {
    if err == nil {
        if ctxErr := MapContextOnly(nil, ctx); ctxErr != nil {
            return ctxErr
        }
        return nil
    }
    if ctxErr := MapContextOnly(err, ctx); ctxErr != nil {
        return ctxErr
    }
    return fmt.Errorf("%s: %w", prefix, err)
}