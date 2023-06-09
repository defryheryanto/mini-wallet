package http

import (
	"io"
	"net/http"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/client"
	"github.com/defryheryanto/mini-wallet/internal/errors"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/request"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/response"
	"github.com/defryheryanto/mini-wallet/internal/wallet"
)

type EnabledWalletResponse struct {
	Id        string    `json:"id"`
	OwnedBy   string    `json:"owned_by"`
	Status    string    `json:"status"`
	EnabledAt time.Time `json:"enabled_at"`
	Balance   float64   `json:"balance"`
}
type DisabledWalletResponse struct {
	Id         string    `json:"id"`
	OwnedBy    string    `json:"owned_by"`
	Status     string    `json:"status"`
	DisabledAt time.Time `json:"disabled_at"`
	Balance    float64   `json:"balance"`
}

type UpdateWalletStatusRequest struct {
	IsDisabled bool `json:"is_disabled"`
}

func HandleEnableWallet(service wallet.WalletIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentClient, err := client.FromContext(r.Context())
		if err != nil {
			response.Failed(w, err)
			return
		}

		targetWallet, err := service.UpdateStatus(r.Context(), currentClient.Xid, true)
		if err != nil {
			response.Failed(w, err)
			return
		}

		response.Success(w, http.StatusCreated, map[string]interface{}{
			"wallet": &EnabledWalletResponse{
				Id:        targetWallet.Id,
				OwnedBy:   targetWallet.OwnedBy,
				EnabledAt: *targetWallet.EnabledAt,
				Status:    targetWallet.Status,
				Balance:   targetWallet.Balance,
			},
		})
	}
}

func HandleViewWallet(service wallet.WalletIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentClient, err := client.FromContext(r.Context())
		if err != nil {
			response.Failed(w, err)
			return
		}

		targetWallet, err := service.GetWalletByXid(r.Context(), currentClient.Xid)
		if err != nil {
			response.Failed(w, err)
			return
		}

		response.Success(w, http.StatusOK, map[string]interface{}{
			"wallet": &EnabledWalletResponse{
				Id:        targetWallet.Id,
				OwnedBy:   targetWallet.OwnedBy,
				EnabledAt: *targetWallet.EnabledAt,
				Status:    targetWallet.Status,
				Balance:   targetWallet.Balance,
			},
		})
	}
}

func HandleUpdateWalletStatus(service wallet.WalletIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := &UpdateWalletStatusRequest{}

		err := request.DecodeBody(r, &requestBody)
		if err != nil {
			if err == io.EOF {
				response.Failed(w, errors.NewValidationError(map[string]interface{}{
					"is_disabled": []string{
						"Missing data for required field.",
					},
				}))
				return
			}
			response.Failed(w, err)
			return
		}

		currentClient, err := client.FromContext(r.Context())
		if err != nil {
			response.Failed(w, err)
			return
		}

		targetWallet, err := service.UpdateStatus(r.Context(), currentClient.Xid, !requestBody.IsDisabled)
		if err != nil {
			response.Failed(w, err)
			return
		}

		if requestBody.IsDisabled {
			response.Success(w, http.StatusOK, map[string]interface{}{
				"wallet": &DisabledWalletResponse{
					Id:         targetWallet.Id,
					OwnedBy:    targetWallet.OwnedBy,
					DisabledAt: *targetWallet.DisabledAt,
					Status:     targetWallet.Status,
					Balance:    targetWallet.Balance,
				},
			})
		} else {
			response.Success(w, http.StatusOK, map[string]interface{}{
				"wallet": &EnabledWalletResponse{
					Id:        targetWallet.Id,
					OwnedBy:   targetWallet.OwnedBy,
					EnabledAt: *targetWallet.EnabledAt,
					Status:    targetWallet.Status,
					Balance:   targetWallet.Balance,
				},
			})
		}

	}
}
