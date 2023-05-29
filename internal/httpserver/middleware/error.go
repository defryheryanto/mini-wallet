package middleware

import "github.com/defryheryanto/mini-wallet/internal/errors"

var ErrInvalidToken = errors.NewUnauthorizedError("authorization token invalid")
