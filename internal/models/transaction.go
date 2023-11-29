package models

import (
	"fanclub/internal/types"
	"time"
)

type Transaction struct {
	ID            int          `json:"id"`
	UserID        int          `json:"user_id"`
	Amount        float64      `json:"amount"`
	Dir           types.TxType `json:"dir"`
	TxHash        string       `json:"tx_hash"`
	WalletAddress string       `json:"wallet_address"`
	CreatedAt     time.Time    `json:"created_at"`
}
