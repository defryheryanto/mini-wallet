package http

import (
	"io"
	"net/http"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/client"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/request"
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

type DepositResponse struct {
	Id          string    `json:"id"`
	DepositedBy string    `json:"deposited_by"`
	Status      string    `json:"status"`
	DepositedAt time.Time `json:"deposited_at"`
	Amount      float64   `json:"amount"`
	ReferenceId string    `json:"reference_id"`
}

type CreateDepositRequest struct {
	Amount      float64 `json:"amount"`
	ReferenceId string  `json:"reference_id"`
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

func HandleCreateDeposit(service transaction.TransactionIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errEmptyReferenceIdMsg := map[string]interface{}{
			"reference_id": []string{
				"Missing data for required field.",
			},
		}
		errEmptyAmountMsg := map[string]interface{}{
			"amount": []string{
				"Missing data for required field.",
			},
		}
		requestBody := &CreateDepositRequest{}

		err := request.DecodeBody(r, &requestBody)
		if err != nil {
			if err == io.EOF {
				response.Failed(w, http.StatusBadRequest, map[string]interface{}{
					"reference_id": errEmptyReferenceIdMsg["reference_id"],
					"amount":       errEmptyAmountMsg["amount"],
				})
				return
			}
			response.Error(w, err)
			return
		}

		if requestBody.ReferenceId == "" {
			response.Failed(w, http.StatusBadRequest, errEmptyReferenceIdMsg)
			return
		}

		currentClient, err := client.FromContext(r.Context())
		if err != nil {
			response.Failed(w, http.StatusUnauthorized, err.Error())
			return
		}

		trx, err := service.CreateDeposit(r.Context(), &transaction.CreateDepositParams{
			CustomerXid: currentClient.Xid,
			ReferenceId: requestBody.ReferenceId,
			Amount:      requestBody.Amount,
		})
		if err != nil {
			if err == transaction.ErrReferenceNoAlreadyExists {
				response.Failed(w, http.StatusBadRequest, err.Error())
				return
			}
			if err == wallet.ErrWalletDisabled {
				response.Failed(w, http.StatusBadRequest, "Wallet disabled")
				return
			}
			if err == wallet.ErrWalletNotFound {
				response.Failed(w, http.StatusNotFound, err.Error())
			}
			response.Error(w, err)
			return
		}

		response.Success(w, http.StatusCreated, &DepositResponse{
			Id:          trx.Id,
			DepositedBy: currentClient.Xid,
			Status:      trx.Status,
			DepositedAt: trx.TransactedAt,
			Amount:      trx.Amount,
			ReferenceId: trx.ReferenceId,
		})
	}
}
