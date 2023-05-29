package wallet

import "fmt"

var ErrOwnedByRequired = fmt.Errorf("owned_by field is required")
var ErrWalletNotFound = fmt.Errorf("wallet not found")
var ErrWalletAlreadyEnabled = fmt.Errorf("wallet already enabled")
var ErrWalletAlreadyDisabled = fmt.Errorf("wallet already disabled")
var ErrWalletDisabled = fmt.Errorf("wallet is not enabled")
var ErrInsufficientBalance = fmt.Errorf("balance insufficient")
