package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/defryheryanto/mini-wallet/internal/app"
	"github.com/go-chi/chi/v5"
)

func HandleRoutes(application *app.Application) http.Handler {
	root := chi.NewRouter()

	root.Get("/api/v1/init", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		json.NewEncoder(w).Encode("Success")
	})

	return root
}
