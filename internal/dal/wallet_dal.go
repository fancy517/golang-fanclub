package dal

import (
	"database/sql"
	"fanclub/internal/models"
)

type WalletDAL interface {
	Insert(models.Wallet) error
	GetOne(int) (models.Wallet, error)
	GetHouseWallet() (float64, error)
	SetUserCredit(userID int, amount float64) error
	AddUserCredit(userID int, amount float64) error
	DeductUserCredit(userID int, amount float64) error
	UpdateLastDepositTime(userID int, lastDeposit int64) error
	GetWalletBalance(userID int) (float64, error)
}

type walletDAL struct {
	DB *sql.DB
}

var _ WalletDAL = (*walletDAL)(nil)

func NewWalletDAL(db *sql.DB) WalletDAL {
	return &walletDAL{db}
}

func (dal *walletDAL) GetWalletBalance(userID int) (float64, error) {
	query := "SELECT credit FROM wallets WHERE user_id=?"
	var credit float64
	err := dal.DB.QueryRow(query, userID).Scan(&credit)
	if err != nil {
		return 0, err
	}
	return credit, nil
}

func (dal *walletDAL) Insert(data models.Wallet) error {
	query := `
	INSERT INTO wallets(user_id, credit, deposit_address, private_key)
	VALUES(?, ?, ?, ?);
	`
	_, err := dal.DB.Exec(query,
		data.UserID,
		data.Credit,
		data.DepositAddress,
		data.PrivateKey,
	)
	return err
}

func (dal *walletDAL) GetOne(userID int) (models.Wallet, error) {
	wallet := models.Wallet{}
	query := `
	SELECT id, user_id, credit, deposit_address, private_key, last_deposit
	FROM wallets
	WHERE user_id = ?;
	`

	err := dal.DB.QueryRow(query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Credit,
		&wallet.DepositAddress,
		&wallet.PrivateKey,
		&wallet.LastDeposit,
	)
	return wallet, err
}

func (dal *walletDAL) SetUserCredit(userID int, amount float64) error {
	query := `
	UPDATE wallets
	SET credit = ?
	WHERE user_id = ?
	`

	_, err := dal.DB.Exec(query, amount, userID)
	return err
}

func (dal *walletDAL) AddUserCredit(userID int, amount float64) error {
	query := `
	UPDATE wallets
	SET credit = ROUND(credit + ?, 3)
	WHERE user_id = ?
	`

	_, err := dal.DB.Exec(query, amount, userID)
	return err
}

func (dal *walletDAL) DeductUserCredit(userID int, amount float64) error {
	return dal.AddUserCredit(userID, -amount)
}

func (dal *walletDAL) UpdateLastDepositTime(userID int, lastDeposit int64) error {
	query := `
	UPDATE wallets
	SET last_deposit = ?
	WHERE user_id = ?;
	`
	_, err := dal.DB.Exec(query, lastDeposit, userID)
	return err
}

func (dal *walletDAL) GetHouseWallet() (float64, error) {
	var totalReward float64 = 0
	var totalWithdraw float64 = 0
	query := `SELECT COALESCE(ROUND(SUM(amount), 5), 0) FROM rewards WHERE user_id = -1;`
	if err := dal.DB.QueryRow(query).Scan(&totalReward); err != nil {
		return 0, err
	}

	query = `SELECT COALESCE(ROUND(SUM(amount), 5), 0) FROM transactions WHERE user_id = -1;`
	if err := dal.DB.QueryRow(query).Scan(&totalWithdraw); err != nil {
		return 0, err
	}

	return totalReward - totalWithdraw, nil
}
