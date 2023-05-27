package gorm

import (
	"time"

	"github.com/defryheryanto/mini-wallet/internal/wallet"
)

type Wallet struct {
	Id         string     `gorm:"primaryKey;column:id"`
	OwnedBy    string     `gorm:"column:owned_by"`
	Status     string     `gorm:"column:status"`
	DisabledAt *time.Time `gorm:"column:disabled_at"`
	EnabledAt  *time.Time `gorm:"column:enabled_at"`
	Balance    float64    `gorm:"column:balance"`
}

func (Wallet) TableName() string {
	return "wallets"
}

func (Wallet) FromServiceModel(data *wallet.Wallet) *Wallet {
	if data == nil {
		return nil
	}

	return &Wallet{
		Id:         data.Id,
		OwnedBy:    data.OwnedBy,
		Status:     data.Status,
		DisabledAt: data.DisabledAt,
		EnabledAt:  data.EnabledAt,
		Balance:    data.Balance,
	}
}

func (w *Wallet) ToServiceModel() *wallet.Wallet {
	return &wallet.Wallet{
		Id:         w.Id,
		OwnedBy:    w.OwnedBy,
		Status:     w.Status,
		DisabledAt: w.DisabledAt,
		EnabledAt:  w.EnabledAt,
		Balance:    w.Balance,
	}
}
