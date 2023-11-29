package types

const (
	ActivateCodeLen         = 10
	ReferralCodeLen         = 9
	MaximumWithdrawalPerDay = 3
)

type DepositTx struct {
	Sender    string
	Amount    float64
	TxHash    string
	Timestamp int64
	// Memo   string
}

type DefaultUri struct {
	ID int `json:"id" uri:"id" binding:"required"`
}
