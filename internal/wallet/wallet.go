package wallet

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	Id         string     `json:"id"`
	OwnedBy    string     `json:"owned_by"`
	Status     string     `json:"status"`
	DisabledAt *time.Time `json:"disabled_at"`
	EnabledAt  *time.Time `json:"enabled_at"`
	Balance    float64    `json:"balance"`
}

type WalletRepository interface {
	Insert(ctx context.Context, data *Wallet) error
	FindById(ctx context.Context, id string) (*Wallet, error)
}

type WalletIService interface {
	Create(ctx context.Context, params *CreateWalletParams) error
}

type WalletService struct {
	repository WalletRepository
}

func NewWalletService(repository WalletRepository) *WalletService {
	return &WalletService{repository}
}

func (s *WalletService) Create(ctx context.Context, params *CreateWalletParams) error {
	if params.OwnedBy == "" {
		return ErrOwnedByRequired
	}

	uuidRandom, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var randomId string
	for {
		randomId = uuidRandom.String()
		existingWallet, err := s.repository.FindById(ctx, randomId)
		if err != nil {
			return err
		}
		if existingWallet == nil {
			break
		}
	}

	err = s.repository.Insert(ctx, &Wallet{
		Id:         randomId,
		OwnedBy:    params.OwnedBy,
		Status:     STATUS_DISABLED,
		DisabledAt: nil,
		EnabledAt:  nil,
		Balance:    0,
	})
	if err != nil {
		return err
	}

	return nil
}
