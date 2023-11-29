package dal

import "database/sql"

type SettingsDAL interface {
	GetLastBlockHeight() (int64, error)
	SetLastBlockHeight(height int64) error
	GetDepositWalletPrivateKey() (string, error)
	SetDepositWalletPrivateKey(key string) error
}

type settingsDAL struct {
	DB *sql.DB
}

var _ SettingsDAL = (*settingsDAL)(nil)

func NewSettingsDAL(db *sql.DB) SettingsDAL {
	return &settingsDAL{db}
}

func (dal *settingsDAL) GetLastBlockHeight() (int64, error) {
	var height int64
	if err := dal.DB.QueryRow(`SELECT last_block_height FROM settings;`).Scan(&height); err != nil {
		return 0, err
	}
	return height, nil
}

func (dal *settingsDAL) SetLastBlockHeight(height int64) error {
	_, err := dal.DB.Exec(`UPDATE settings SET last_block_height = ?;`, height)
	return err
}

func (dal *settingsDAL) GetDepositWalletPrivateKey() (string, error) {
	var key string
	if err := dal.DB.QueryRow(`SELECT deposit_wallet_private_key FROM settings;`).Scan(&key); err != nil {
		return "", err
	}
	return key, nil
}

func (dal *settingsDAL) SetDepositWalletPrivateKey(key string) error {
	_, err := dal.DB.Exec(`UPDATE settings SET deposit_wallet_private_key = ?;`, key)
	return err
}
