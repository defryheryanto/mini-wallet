package wallet

import (
	"github.com/defryheryanto/mini-wallet/internal/errors"
)

var ErrOwnedByRequired = errors.NewValidationError("owned_by field is required")
var ErrWalletNotFound = errors.NewNotFoundError("wallet not found")
var ErrWalletAlreadyEnabled = errors.NewValidationError("Already enabled")
var ErrWalletAlreadyDisabled = errors.NewValidationError("Already disabled")
var ErrWalletDisabled = errors.NewNotFoundError("Wallet disabled")
var ErrInsufficientBalance = errors.NewValidationError("balance insufficient")
