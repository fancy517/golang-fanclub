package troner

import (
	"fanclub/internal/interfaces"
	tronclient "fanclub/pkg/go-tron/client"
)

type troner struct {
	Host            string
	AddrUSDT        string
	HouseAddress    string
	HousePrivateKey string
	MinDeposit      float64
	EarningFee      float64
}

var _ interfaces.Troner = (*troner)(nil)

func NewTroner(host, houseAddr, houseKey string, minAmount float64, earningFee float64) interfaces.Troner {
	return &troner{
		Host:            host,
		HouseAddress:    houseAddr,
		HousePrivateKey: houseKey,
		MinDeposit:      minAmount,
		EarningFee:      earningFee,
	}
}

func (t *troner) NewClient() *tronclient.Client {
	return tronclient.New(t.Host)
}

func (t *troner) SetUSDTContractAddress(addr string) {
	t.AddrUSDT = addr
}

func (t *troner) GetMinDepositValue() float64 {
	return t.MinDeposit
}

func (t *troner) GetLastBlockHeight() (int64, error) {
	client := t.NewClient()
	block, err := client.GetLatestBlock()
	if err != nil {
		return 0, err
	}
	return int64(block.BlockHeader.RawData.Number), nil
}

func (t *troner) GetAccountBalance(addr string) (float64, error) {
	client := t.NewClient()
	acc, err := client.GetAccount(addr)
	if err != nil {
		return 0, err
	}
	return float64(acc.Balance) / float64(1000000), nil
}
