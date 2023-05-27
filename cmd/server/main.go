package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/defryheryanto/mini-wallet/internal/httpserver"
)

func main() {
	var appServer *http.Server

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		gormClient := setupGormClient()
		appContainer := buildApp(gormClient)
		appServer = &http.Server{
			Addr:    ":8080",
			Handler: httpserver.HandleRoutes(appContainer),
		}

		log.Println("starting server on port 8080")
		if err := appServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error starting server: %v\n", err)
		}
	}()

	<-done

	log.Println("shutting down server")
	if err := appServer.Shutdown(ctx); err != nil {
		log.Printf("error shutting down server: %v\n", err)
	}
	cancel()

	log.Println("server shutdown gracefully")
}
