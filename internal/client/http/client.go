package http

import (
	"io"
	"net/http"

	"github.com/defryheryanto/mini-wallet/internal/client"
	"github.com/defryheryanto/mini-wallet/internal/errors"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/request"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/response"
)

type CreateClientRequest struct {
	CustomerXid string `json:"customer_xid"`
}

type CreateClientResponse struct {
	Token string `json:"token"`
}

func HandleCreateClient(service client.ClientIService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errEmptyCustomerXidData := map[string]interface{}{
			"customer_xid": []string{
				"Missing data for required field.",
			},
		}
		requestBody := &CreateClientRequest{}

		err := request.DecodeBody(r, &requestBody)
		if err != nil {
			if err == io.EOF {
				response.Failed(w, errors.NewValidationError(errEmptyCustomerXidData))
				return
			}
			response.Failed(w, err)
			return
		}

		if requestBody.CustomerXid == "" {
			response.Failed(w, errors.NewValidationError(errEmptyCustomerXidData))
			return
		}

		createdClient, err := service.Create(r.Context(), requestBody.CustomerXid)
		if err != nil {
			if err == client.ErrXidAlreadyTaken {
				response.Failed(w, err)
				return
			}
			response.Failed(w, err)
			return
		}

		response.Success(w, http.StatusCreated, &CreateClientResponse{
			Token: createdClient.Token,
		})
	}
}
