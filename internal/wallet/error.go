package wallet

import "fmt"

var ErrOwnedByRequired = fmt.Errorf("owned_by field is required")
var ErrWalletNotFound = fmt.Errorf("wallet not found")
var ErrWalletAlreadyEnabled = fmt.Errorf("wallet already enabled")
var ErrWalletDisabled = fmt.Errorf("wallet is not enabled")
