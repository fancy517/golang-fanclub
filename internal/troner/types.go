package troner

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance to collect")
	ErrBelowMinDeposit     = errors.New("amount is below min deposit")
)

const (
	Mainnet       = "https://api.trongrid.io"
	ShastaTestnet = "https://api.shasta.trongrid.io"
	NileTestnet   = "https://nile.trongrid.io"
)
