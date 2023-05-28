package main

import (
	"github.com/defryheryanto/mini-wallet/internal/app"
	"github.com/defryheryanto/mini-wallet/internal/client"
	client_repository "github.com/defryheryanto/mini-wallet/internal/client/repository/gorm"
	"github.com/defryheryanto/mini-wallet/internal/storage/manager"
	gorm_storage_manager "github.com/defryheryanto/mini-wallet/internal/storage/manager/gorm"
	"github.com/defryheryanto/mini-wallet/internal/transaction"
	transaction_repository "github.com/defryheryanto/mini-wallet/internal/transaction/repository/gorm"
	"github.com/defryheryanto/mini-wallet/internal/wallet"
	wallet_repository "github.com/defryheryanto/mini-wallet/internal/wallet/repository/gorm"
	"gorm.io/gorm"
)

func buildApp(db *gorm.DB) *app.Application {
	gormManager := setupGormStorageManager(db)
	walletService := setupWallet(db)
	clientService := setupClient(db, walletService, gormManager)
	transactionService := setupTransaction(db, walletService)

	return &app.Application{
		WalletService:      walletService,
		ClientService:      clientService,
		TransactionService: transactionService,
	}
}

func setupGormStorageManager(db *gorm.DB) manager.StorageManager {
	return gorm_storage_manager.NewGormStorageManager(db)
}

func setupWallet(db *gorm.DB) wallet.WalletIService {
	repository := wallet_repository.NewWalletRepository(db)
	return wallet.NewWalletService(repository)
}

func setupClient(db *gorm.DB, walletService wallet.WalletIService, storageManager manager.StorageManager) client.ClientIService {
	repository := client_repository.NewClientRepository(db)
	return client.NewClientService(repository, walletService, storageManager)
}

func setupTransaction(db *gorm.DB, walletService wallet.WalletIService) transaction.TransactionIService {
	repository := transaction_repository.NewTransactionRepository(db)
	return transaction.NewTransactionService(repository, walletService)
}
