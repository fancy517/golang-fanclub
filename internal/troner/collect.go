package troner

import (
	tronacc "fanclub/pkg/go-tron/account"
	"fmt"
)

func (t *troner) GetEarningFee() float64 {
	return t.EarningFee
}

func (t *troner) Collect2House(senderPrivateKey string, amount float64) (string, error) {
	if amount < t.MinDeposit {
		return "", ErrBelowMinDeposit
	}

	sender, err := tronacc.FromPrivateKeyHex(senderPrivateKey)
	if err != nil {
		return "", err
	}

	senderAddress := sender.Address().ToBase58()
	balance, err := t.GetAccountBalance(senderAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get account balance, address=%v; %w", senderAddress, err)
	}

	fmt.Printf("balance: %v, amount: %v\n", balance, amount)
	if balance < amount {
		return "", ErrInsufficientBalance
	}

	resource, err := t.NewClient().GetAccountResource(senderAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get account resource, address=%v; %w", senderAddress, err)
	}

	tx, err := t.GetTransaction(senderPrivateKey, t.HouseAddress, Trx2Sun(amount))
	if err != nil {
		return "", fmt.Errorf("failed to get tx; %w", err)
	}

	fee := len(*tx.RawDataHex)
	var txId string
	if fee <= resource.FreeNetAvailable() {
		txId, err = t.TransferTRX(senderPrivateKey, t.HouseAddress, Trx2Sun(amount))
		if err != nil {
			return "", err
		}
	} else {
		txId, err = t.TransferTRX(senderPrivateKey, t.HouseAddress, Trx2Sun(amount-TxSize2Fee(fee)))
		if err != nil {
			return "", err
		}
	}

	return txId, nil
}
