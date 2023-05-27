package httpserver

import (
	"net/http"

	"github.com/defryheryanto/mini-wallet/internal/app"
	client_http "github.com/defryheryanto/mini-wallet/internal/client/http"
	"github.com/go-chi/chi/v5"
)

func HandleRoutes(application *app.Application) http.Handler {
	root := chi.NewRouter()

	root.Post("/api/v1/init", client_http.HandleCreateClient(application.ClientService))

	return root
}
