package wallet

type WithdrawRequest struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

type WithdrawResponse struct {
	TxHash  string  `json:"tx_hash"`
	Balance float64 `json:"balance"`
}

type GetBalanceResponse struct {
	Balance        float64 `json:"balance"`
	DepositAddress string  `json:"deposit_address"`
}

type CheckDepositTxsResponse struct {
	Balance    float64 `json:"balance"`
	NewDeposit bool    `json:"new_deposit"`
}
