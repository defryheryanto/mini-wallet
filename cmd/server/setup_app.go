package main

import (
	"github.com/defryheryanto/mini-wallet/internal/app"
	"github.com/defryheryanto/mini-wallet/internal/client"
	client_repository "github.com/defryheryanto/mini-wallet/internal/client/repository/gorm"
	"github.com/defryheryanto/mini-wallet/internal/wallet"
	wallet_repository "github.com/defryheryanto/mini-wallet/internal/wallet/repository/gorm"
	"gorm.io/gorm"
)

func buildApp(db *gorm.DB) *app.Application {
	walletService := setupWallet(db)
	clientService := setupClient(db, walletService)

	return &app.Application{
		WalletService: walletService,
		ClientService: clientService,
	}
}

func setupWallet(db *gorm.DB) wallet.WalletIService {
	repository := wallet_repository.NewWalletRepository(db)
	return wallet.NewWalletService(repository)
}

func setupClient(db *gorm.DB, walletService wallet.WalletIService) client.ClientIService {
	repository := client_repository.NewClientRepository(db)
	return client.NewClientService(repository, walletService)
}
