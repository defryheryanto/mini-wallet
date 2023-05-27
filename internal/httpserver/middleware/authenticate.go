package middleware

import (
	"net/http"
	"strings"

	"github.com/defryheryanto/mini-wallet/internal/client"
	"github.com/defryheryanto/mini-wallet/internal/httpserver/response"
)

func AuthenticateClient(clientService client.ClientIService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := strings.Split(r.Header.Get("Authorization"), "Token ")
			if len(authorization) < 2 {
				response.Failed(w, http.StatusUnauthorized, "authorization token invalid")
				return
			}

			token := authorization[1]

			currentClient, err := clientService.GetByToken(r.Context(), token)
			if err != nil {
				response.Error(w, err)
				return
			}
			if currentClient == nil {
				response.Failed(w, http.StatusUnauthorized, "authorization token invalid")
				return
			}

			ctx := client.Inject(r.Context(), currentClient)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
