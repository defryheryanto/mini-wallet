package http

import (
	"net/http"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/client"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/response"
	"github.com/defryheryanto/mini-wallet/internal/transaction"
	"github.com/defryheryanto/mini-wallet/internal/wallet"
)

type TransactionResponse struct {
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	TransactedAt time.Time `json:"transacted_at"`
	Type         string    `json:"type"`
	Amount       float64   `json:"amount"`
	ReferenceId  string    `json:"reference_id"`
}

func HandleGetWalletTransactions(service transaction.TransactionIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentClient, err := client.FromContext(r.Context())
		if err != nil {
			response.Failed(w, http.StatusUnauthorized, err.Error())
			return
		}

		transactions, err := service.GetTransactionsByCustomerXid(r.Context(), currentClient.Xid)
		if err != nil {
			if err == wallet.ErrWalletNotFound {
				response.Failed(w, http.StatusNotFound, err.Error())
				return
			}
			if err == wallet.ErrWalletDisabled {
				response.Failed(w, http.StatusBadRequest, "Wallet disabled")
				return
			}
			response.Error(w, err)
			return
		}

		trx := []*TransactionResponse{}

		for _, tr := range transactions {
			trx = append(trx, &TransactionResponse{
				Id:           tr.Id,
				Status:       tr.Status,
				TransactedAt: tr.TransactedAt,
				Type:         tr.Type,
				Amount:       tr.Amount,
				ReferenceId:  tr.ReferenceId,
			})
		}

		response.Success(w, http.StatusOK, map[string]interface{}{
			"transactions": trx,
		})
	}
}
