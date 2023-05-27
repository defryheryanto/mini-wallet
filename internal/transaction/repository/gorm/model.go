package gorm

import (
	"time"

	"github.com/defryheryanto/mini-wallet/internal/transaction"
)

type Transaction struct {
	Id           string    `gorm:"primaryKey;column:id"`
	Status       string    `gorm:"column:status"`
	TransactedAt time.Time `gorm:"column:transacted_at"`
	Type         string    `gorm:"column:type"`
	Amount       float64   `gorm:"column:amount"`
	ReferenceId  string    `gorm:"column:reference_id"`
	WalletId     string    `gorm:"column:wallet_id"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (Transaction) FromServiceModel(data *transaction.Transaction) *Transaction {
	if data == nil {
		return nil
	}

	return &Transaction{
		Id:           data.Id,
		Status:       data.Status,
		TransactedAt: data.TransactedAt,
		Type:         data.Type,
		Amount:       data.Amount,
		ReferenceId:  data.ReferenceId,
		WalletId:     data.WalletId,
	}
}

func (c *Transaction) ToServiceModel() *transaction.Transaction {
	return &transaction.Transaction{
		Id:           c.Id,
		Status:       c.Status,
		TransactedAt: c.TransactedAt,
		Type:         c.Type,
		Amount:       c.Amount,
		ReferenceId:  c.ReferenceId,
		WalletId:     c.WalletId,
	}
}

func SliceToServiceModel(data []*Transaction) []*transaction.Transaction {
	if data == nil {
		return nil
	}

	transactions := []*transaction.Transaction{}
	for _, trx := range data {
		transactions = append(transactions, trx.ToServiceModel())
	}

	return transactions
}
