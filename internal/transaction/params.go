package transaction

type CreateDepositParams struct {
	CustomerXid string  `json:"customer_xid"`
	ReferenceId string  `json:"reference_no"`
	Amount      float64 `json:"amount"`
}

type CreateWithdrawalParams struct {
	CustomerXid string  `json:"customer_xid"`
	ReferenceId string  `json:"reference_no"`
	Amount      float64 `json:"amount"`
}
