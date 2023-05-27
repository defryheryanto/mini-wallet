package http

import (
	"net/http"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/client"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/response"
	"github.com/defryheryanto/mini-wallet/internal/wallet"
)

type EnableWalletResponse struct {
	Id        string    `json:"id"`
	OwnedBy   string    `json:"owned_by"`
	Status    string    `json:"status"`
	EnabledAt time.Time `json:"enabled_at"`
	Balance   float64   `json:"balance"`
}

func HandleEnableWallet(service wallet.WalletIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentClient, err := client.FromContext(r.Context())
		if err != nil {
			response.Failed(w, http.StatusUnauthorized, err.Error())
			return
		}

		targetWallet, err := service.EnableWallet(r.Context(), currentClient.Xid)
		if err != nil {
			if err == wallet.ErrWalletNotFound {
				response.Failed(w, http.StatusNotFound, err.Error())
				return
			}
			if err == wallet.ErrWalletAlreadyEnabled {
				response.Failed(w, http.StatusBadRequest, "Already enabled")
				return
			}
			response.Error(w, err)
			return
		}

		response.Success(w, http.StatusCreated, map[string]interface{}{
			"wallet": &EnableWalletResponse{
				Id:        targetWallet.Id,
				OwnedBy:   targetWallet.OwnedBy,
				EnabledAt: *targetWallet.EnabledAt,
				Status:    targetWallet.Status,
				Balance:   targetWallet.Balance,
			},
		})
	}
}