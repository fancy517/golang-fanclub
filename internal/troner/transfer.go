package troner

import (
	"errors"
	"fanclub/pkg/go-tron"
	tronabi "fanclub/pkg/go-tron/abi"
	tronacc "fanclub/pkg/go-tron/account"
	tronaddr "fanclub/pkg/go-tron/address"
	tronclient "fanclub/pkg/go-tron/client"
	"fmt"
	"os"
)

// Returns TxID and error value
func (t *troner) TransferTRX(senderPrivateKey, receiverAddress string, amount uint64) (string, error) {
	client := tronclient.New(t.Host)

	tx, err := t.GetTransaction(senderPrivateKey, receiverAddress, amount)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction, sender=%v, receiver=%v, amount=%v; %w",
			senderPrivateKey, receiverAddress, amount, err)
	}

	if err := client.BroadcastTransaction(&tx); err != nil {
		return "", err
	}

	return tx.Id, nil
}

func (t *troner) GetTransaction(senderPrivateKey, receiverAddress string, amount uint64) (tron.Transaction, error) {
	client := t.NewClient()
	sender, err := tronacc.FromPrivateKeyHex(senderPrivateKey)
	if err != nil {
		return tron.Transaction{}, err
	}

	receiver, err := tronaddr.FromBase58(receiverAddress)
	if err != nil {
		return tron.Transaction{}, err
	}

	return client.Transfer(sender, receiver, amount)
}

func (t *troner) Withdraw(address string, amount float64) (string, error) {
	return t.TransferTRX(t.HousePrivateKey, address, Trx2Sun(amount))
}

// Returns TxID and error value
func (t *troner) TransferUSDT(senderPrivateKey, receiverAddress string, amount uint64) (string, error) {
	if t.AddrUSDT == "" {
		return "", errors.New("usdt contract address is not valid")
	}

	abiData, err := os.ReadFile("./abi/trc-20.json")
	if err != nil {
		return "", err
	}

	abi := tronabi.ABI{}
	if err := abi.UnmarshalJSON(abiData); err != nil {
		return "", err
	}

	client := tronclient.New(t.Host)
	sender, err := tronacc.FromPrivateKeyHex(senderPrivateKey)
	if err != nil {
		return "", err
	}

	contractAddr, err := tronaddr.FromBase58(t.AddrUSDT)
	if err != nil {
		return "", err
	}

	receiver, err := tronaddr.FromBase58(receiverAddress)
	if err != nil {
		return "", err
	}

	input := tronclient.CallContractInput{
		Address:   contractAddr,
		Function:  abi.Functions["transfer"],
		Arguments: []interface{}{receiver, amount},
		FeeLimit:  1000000000,
		CallValue: 0,
	}

	tx, result, err := client.TriggerSmartContract(sender, input)
	if err != nil {
		return "", err
	}

	if !result {
		return "", errors.New("failed to get tx")
	}

	if sender.Sign(tx) != nil {
		return "", errors.New("failed to sign tx")
	}

	if err := client.BroadcastTransaction(tx); err != nil {
		return "", err
	}

	return tx.Id, nil
}
