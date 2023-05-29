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
	FindByCustomerXid(ctx context.Context, xid string) (*Wallet, error)
	Update(ctx context.Context, data *Wallet) error
}

type WalletIService interface {
	Create(ctx context.Context, params *CreateWalletParams) error
	UpdateStatus(ctx context.Context, customerXid string, isEnabled bool) (*Wallet, error)
	GetWalletByXid(ctx context.Context, customerXid string) (*Wallet, error)
	AddBalance(ctx context.Context, walletId string, amount float64) error
	ValidateWallet(target *Wallet) error
	DeductBalance(ctx context.Context, walletId string, amount float64) error
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

func (s *WalletService) UpdateStatus(ctx context.Context, customerXid string, isEnabled bool) (*Wallet, error) {
	currentWallet, err := s.repository.FindByCustomerXid(ctx, customerXid)
	if err != nil {
		return nil, err
	}
	if currentWallet == nil {
		return nil, ErrWalletNotFound
	}

	now := time.Now()
	if isEnabled {
		if currentWallet.Status == STATUS_ENABLED {
			return nil, ErrWalletAlreadyEnabled
		}
		currentWallet.DisabledAt = nil
		currentWallet.EnabledAt = &now
		currentWallet.Status = STATUS_ENABLED
	} else {
		if currentWallet.Status == STATUS_DISABLED {
			return nil, ErrWalletAlreadyDisabled
		}
		currentWallet.DisabledAt = &now
		currentWallet.EnabledAt = nil
		currentWallet.Status = STATUS_DISABLED
	}

	err = s.repository.Update(ctx, currentWallet)
	if err != nil {
		return nil, err
	}

	currentWallet, err = s.repository.FindByCustomerXid(ctx, customerXid)
	if err != nil {
		return nil, err
	}

	return currentWallet, nil
}

func (s *WalletService) GetWalletByXid(ctx context.Context, customerXid string) (*Wallet, error) {
	currentWallet, err := s.repository.FindByCustomerXid(ctx, customerXid)
	if err != nil {
		return nil, err
	}
	if currentWallet == nil {
		return nil, ErrWalletNotFound
	}
	if currentWallet.Status != STATUS_ENABLED {
		return nil, ErrWalletDisabled
	}

	return currentWallet, nil
}

func (s *WalletService) AddBalance(ctx context.Context, walletId string, amount float64) error {
	// TODO: This process is sensitive to racing condition
	// TODO: Implement redis lock here to avoid it
	targetWallet, err := s.repository.FindById(ctx, walletId)
	if err != nil {
		return err
	}
	if err = s.ValidateWallet(targetWallet); err != nil {
		return err
	}

	targetWallet.Balance += amount

	err = s.repository.Update(ctx, targetWallet)
	if err != nil {
		return err
	}

	return nil
}

func (s *WalletService) ValidateWallet(target *Wallet) error {
	if target == nil {
		return ErrWalletNotFound
	}
	if target.Status == STATUS_DISABLED {
		return ErrWalletDisabled
	}

	return nil
}

func (s *WalletService) DeductBalance(ctx context.Context, walletId string, amount float64) error {
	// TODO: This process is sensitive to racing condition
	// TODO: Implement redis lock here to avoid it
	targetWallet, err := s.repository.FindById(ctx, walletId)
	if err != nil {
		return err
	}
	if err := s.ValidateWallet(targetWallet); err != nil {
		return err
	}

	if targetWallet.Balance < amount {
		return ErrInsufficientBalance
	}

	targetWallet.Balance -= amount
	err = s.repository.Update(ctx, targetWallet)
	if err != nil {
		return err
	}

	return nil
}
