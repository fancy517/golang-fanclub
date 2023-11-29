package dal

import "database/sql"

type CryptoPriceDAL interface {
	Update(token string, usdValue float64) error
	GetUsdPrice(token string) (float64, error)
}

type cryptopriceDAL struct {
	DB *sql.DB
}

var _ CryptoPriceDAL = (*cryptopriceDAL)(nil)

func NewCryptoPriceDAL(db *sql.DB) CryptoPriceDAL {
	return &cryptopriceDAL{db}
}

func (dal *cryptopriceDAL) Update(token string, usdValue float64) error {
	query := `REPLACE INTO crypto_prices (token, usd, updated_at) VALUES (?, ?, NOW());`
	_, err := dal.DB.Exec(query, token, usdValue)
	return err
}

func (dal *cryptopriceDAL) GetUsdPrice(token string) (float64, error) {
	var price float64 = 0
	query := `SELECT usd FROM crypto_prices WHERE token = ?;`
	err := dal.DB.QueryRow(query, token).Scan(&price)
	return price, err
}
