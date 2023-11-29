package crypto

import (
	"encoding/json"
	"fanclub/internal/dal"
	"fanclub/internal/interfaces"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const CollectInterval = time.Minute * 2
const CoingeckoUrl = "https://api.coingecko.com/api/v3/simple/price?ids=tron&vs_currencies=usd"

type priceCollector struct {
	dal dal.AppDAL
}

func NewCryptoPriceCollector(dal dal.AppDAL) interfaces.CryptoPriceCollector {
	return &priceCollector{dal}
}

func (c *priceCollector) Run() {
	go func() {
		value, err := c.fetchPriceFromCoingecko()
		if err != nil {
			log.Println(fmt.Errorf("failed to collect tron price; %w", err))
		} else {
			if err := c.updatePrice(value); err != nil {
				log.Println(fmt.Errorf("failed to update tron price; %w", err))
			}
		}

		time.Sleep(CollectInterval)
	}()
}

func (c *priceCollector) fetchPriceFromCoingecko() (float64, error) {
	resp, err := http.Get(CoingeckoUrl)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil
	}

	var req CoingeckoResponse
	err = json.Unmarshal(bs, &req)
	return req.Tron.USD, err
}

func (c *priceCollector) updatePrice(tronUsdPrice float64) error {
	return c.dal.CryptoPrice.Update("trx", tronUsdPrice)
}
