package interfaces

import "fanclub/internal/types"

type Troner interface {
	GetMinDepositValue() float64
	GetLastBlockHeight() (int64, error)
	GetAccountBalance(addr string) (float64, error)
	TransferTRX(senderPrivateKey, receiverAddress string, amount uint64) (string, error)
	GetDepositTxs(depositAddr string, lastDeposit int64) ([]types.DepositTx, error)
	Collect2House(senderPrivateKey string, amount float64) (string, error)
	Withdraw(address string, amount float64) (string, error)
	GetEarningFee() float64
}

type CryptoPriceCollector interface {
	Run()
}

type App interface {
	Serve()
	CloseDB()
}
