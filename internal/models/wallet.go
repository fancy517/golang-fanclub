package models

import "time"

type Wallet struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Credit         float64   `json:"credit"`
	DepositAddress string    `json:"deposit_address"`
	PrivateKey     string    `json:"private_key"`
	LastDeposit    int64     `json:"last_deposit"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
