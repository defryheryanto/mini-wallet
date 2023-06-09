package httpserver

import (
	"net/http"

	"github.com/defryheryanto/mini-wallet/internal/app"
	client_http "github.com/defryheryanto/mini-wallet/internal/client/http"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/middleware"
	transaction_http "github.com/defryheryanto/mini-wallet/internal/transaction/http"
	wallet_http "github.com/defryheryanto/mini-wallet/internal/wallet/http"
	"github.com/go-chi/chi/v5"
)

func HandleRoutes(application *app.Application) http.Handler {
	root := chi.NewRouter()

	root.Post("/api/v1/init", client_http.HandleCreateClient(application.ClientService))

	root.Group(func(r chi.Router) {
		r.Use(middleware.AuthenticateClient(application.ClientService))

		r.Get("/api/v1/wallet", wallet_http.HandleViewWallet(application.WalletService))
		r.Post("/api/v1/wallet", wallet_http.HandleEnableWallet(application.WalletService))
		r.Patch("/api/v1/wallet", wallet_http.HandleUpdateWalletStatus(application.WalletService))

		r.Get("/api/v1/wallet/transactions", transaction_http.HandleGetWalletTransactions(application.TransactionService))
		r.Post("/api/v1/wallet/deposits", transaction_http.HandleCreateDeposit(application.TransactionService))
		r.Post("/api/v1/wallet/withdrawals", transaction_http.HandleCreateWithdrawal(application.TransactionService))
	})

	return root
}
