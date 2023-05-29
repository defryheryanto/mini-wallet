package transaction

import (
	"github.com/defryheryanto/mini-wallet/internal/errors"
)

var ErrReferenceNoAlreadyExists = errors.NewValidationError("reference number already exists")
var ErrEmptyCustomerXid = errors.NewValidationError("customer xid is required")
var ErrEmptyReferenceId = errors.NewValidationError("reference id is required")
