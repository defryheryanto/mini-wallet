package app

import (
	"github.com/defryheryanto/mini-wallet/internal/client"
	"github.com/defryheryanto/mini-wallet/internal/transaction"
	"github.com/defryheryanto/mini-wallet/internal/wallet"
)

type Application struct {
	WalletService      wallet.WalletIService
	ClientService      client.ClientIService
	TransactionService transaction.TransactionIService
}
