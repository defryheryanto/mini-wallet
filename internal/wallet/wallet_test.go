package wallet_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/defryheryanto/mini-wallet/internal/wallet"
	"github.com/defryheryanto/mini-wallet/internal/wallet/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletService_Create(t *testing.T) {
	mockedErr := fmt.Errorf("mocked")

	t.Run("should return error when params invalid", func(t *testing.T) {
		service := wallet.NewWalletService(mocks.NewWalletRepository(t))

		err := service.Create(context.TODO(), &wallet.CreateWalletParams{})
		assert.Equal(t, wallet.ErrOwnedByRequired, err)
	})

	t.Run("should return error when failed to get wallet by id", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, mockedErr)
		service := wallet.NewWalletService(repository)

		err := service.Create(context.TODO(), &wallet.CreateWalletParams{
			OwnedBy: "test",
		})
		assert.Equal(t, mockedErr, err)
	})

	t.Run("should return error when failed to insert", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(mockedErr)
		service := wallet.NewWalletService(repository)

		err := service.Create(context.TODO(), &wallet.CreateWalletParams{
			OwnedBy: "test",
		})
		assert.Equal(t, mockedErr, err)
	})

	t.Run("should not raise error if operation success", func(t *testing.T) {
		repository := mocks.NewWalletRepository(t)
		repository.On("FindById", mock.Anything, mock.Anything).Return(nil, nil)
		repository.On("Insert", mock.Anything, mock.Anything).Return(nil)

		service := wallet.NewWalletService(repository)
		err := service.Create(context.TODO(), &wallet.CreateWalletParams{
			OwnedBy: "test",
		})
		assert.Nil(t, err)
	})
}
