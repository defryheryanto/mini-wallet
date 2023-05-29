package client

import (
	"github.com/defryheryanto/mini-wallet/internal/errors"
)

var ErrXidAlreadyTaken = errors.NewValidationError("xid already taken")
var ErrInvalidClient = errors.NewUnauthorizedError("client invalid")
