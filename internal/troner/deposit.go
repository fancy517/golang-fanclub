package troner

import (
	"encoding/json"
	"errors"
	"fanclub/internal/types"
	"fmt"
	"io"
	"net/http"
)

type TxRawData struct {
	Data     string `json:"data,omitempty"`
	Contract []struct {
		Parameter struct {
			Value struct {
				Amount       int64  `json:"amount"`
				OwnerAddress string `json:"owner_address"`
				ToAddress    string `json:"to_address"`
			} `json:"value"`
		} `json:"parameter"`
		Type string `json:"type"`
	} `json:"contract"`
}

type AccountTxResponse struct {
	Data []struct {
		TxID           string    `json:"txID"`
		BlockNumber    int64     `json:"blockNumber"`
		BlockTimestamp int64     `json:"block_timestamp"`
		RawData        TxRawData `json:"raw_data"`
	} `json:"data"`
	Success bool `json:"success"`
	Meta    struct {
		At       int64 `json:"at"`
		PageSize int   `json:"page_size"`
	} `json:"meta"`
}

func (t *troner) GetDepositTxs(depositAddr string, lastDeposit int64) ([]types.DepositTx, error) {
	accountTxURL := fmt.Sprintf("%s/v1/accounts/%s/transactions?only_to=true&min_timestamp=%v", t.Host, depositAddr, lastDeposit)
	response, err := http.Get(accountTxURL)
	if err != nil {
		return nil, fmt.Errorf("failed to account incoming tx; %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("account tx response read err; %w", err)
	}

	var txResponse AccountTxResponse
	if err := json.Unmarshal(body, &txResponse); err != nil {
		return nil, fmt.Errorf("account tx response json unmarshal error; %w", err)
	}

	if !txResponse.Success {
		return nil, errors.New("account tx response success is false")
	}

	deposits := make([]types.DepositTx, 0)
	for _, tx := range txResponse.Data {
		if len(tx.RawData.Contract) > 0 {
			dtx := types.DepositTx{
				Sender:    tx.RawData.Contract[0].Parameter.Value.OwnerAddress,
				Amount:    float64(tx.RawData.Contract[0].Parameter.Value.Amount) / float64(1000000),
				TxHash:    tx.TxID,
				Timestamp: tx.BlockTimestamp,
			}
			deposits = append(deposits, dtx)
		}
	}

	return deposits, nil
}
