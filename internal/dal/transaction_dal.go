package dal

import (
	"database/sql"
	"fanclub/internal/models"
	"fanclub/internal/types"
	"fmt"
)

type TransactionDAL interface {
	AddDepositHistory(userID int, txhash, address string, amount float64) error
	AddWithdrawalHistory(userID int, txhash, address string, amount float64) error
	AddCollectHistory(userID int, txhash, address string, amount float64) error

	GetUserTxs(userID int, filter types.TxFilter) ([]models.Transaction, int, error)
	GetTodaysWithdrawlCount(userID int) (int, error)
	GetTransactions(userID int) ([]models.Transaction, error)
}

type transactionDAL struct {
	DB *sql.DB
}

var _ TransactionDAL = (*transactionDAL)(nil)

func NewTransactionDAL(db *sql.DB) TransactionDAL {
	return &transactionDAL{db}
}

func (dal *transactionDAL) GetTransactions(userID int) ([]models.Transaction, error) {
	query := "SELECT id, user_id, amount, dir, tx_hash, wallet_address, created_at FROM transactions WHERE user_id = ? AND dir != ?"
	rows, err := dal.DB.Query(query, userID, "C")
	if err != nil {
		fmt.Printf("error :%v\n", err)
		return nil, err
	}
	defer rows.Close()
	res := make([]models.Transaction, 0)
	for rows.Next() {
		_txs := models.Transaction{}
		err := rows.Scan(
			&_txs.ID,
			&_txs.UserID,
			&_txs.Amount,
			&_txs.Dir,
			&_txs.TxHash,
			&_txs.WalletAddress,
			&_txs.CreatedAt,
		)
		if err != nil {
			fmt.Printf("error :%v\n", err)
			return nil, err
		}
		res = append(res, _txs)
	}
	return res, nil
}

func (dal *transactionDAL) AddHistory(userID int, txhash, address string, amount float64, dir types.TxType) error {
	query := `
	INSERT INTO transactions (user_id, amount, dir, tx_hash, wallet_address)
	VALUES (?, ?, ?, ?, ?);
	`
	_, err := dal.DB.Exec(query, userID, amount, dir, txhash, address)
	return err
}

func (dal *transactionDAL) AddWithdrawalHistory(userID int, txhash, address string, amount float64) error {
	return dal.AddHistory(userID, txhash, address, amount, types.TxTypeWithdrawal)
}

func (dal *transactionDAL) AddDepositHistory(userID int, txhash, address string, amount float64) error {
	return dal.AddHistory(userID, txhash, address, amount, types.TxTypeDeposit)
}

func (dal *transactionDAL) AddCollectHistory(userID int, txhash, address string, amount float64) error {
	return dal.AddHistory(userID, txhash, address, amount, types.TxTypeCollect)
}

func (dal *transactionDAL) GetUserTxs(userID int, filter types.TxFilter) ([]models.Transaction, int, error) {
	typeFilter := ""
	if filter.Type == "Deposit" {
		typeFilter = fmt.Sprintf("AND dir = '%s'", types.TxTypeDeposit)
	} else if filter.Type == "Withdrawal" {
		typeFilter = fmt.Sprintf("AND dir = '%s'", types.TxTypeWithdrawal)
	}

	dateFilter := ""
	if filter.Year != 0 {
		dateFilter = fmt.Sprintf("AND YEAR(created_at) = %d AND MONTH(created_at) = %d", filter.Year, filter.Month)
	}

	query := `
	SELECT COUNT(*)
	FROM transactions t
	WHERE user_id = ?
		AND dir <> 'C'
		` + typeFilter + `
		` + dateFilter

	total := 0
	if err := dal.DB.QueryRow(query, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query = `
	SELECT *
	FROM transactions t
	WHERE user_id = ?
		AND dir <> 'C'
		` + typeFilter + `
		` + dateFilter + `
	ORDER BY id desc
	LIMIT ? OFFSET ?;
	`
	rows, err := dal.DB.Query(query, userID, filter.PageSize, filter.PageSize*(filter.Page-1))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	txs := make([]models.Transaction, 0)
	for rows.Next() {
		tx := models.Transaction{}
		if err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Amount,
			&tx.Dir,
			&tx.TxHash,
			&tx.WalletAddress,
			&tx.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		txs = append(txs, tx)
	}
	return txs, total, nil
}

func (dal *transactionDAL) GetTodaysWithdrawlCount(userID int) (count int, err error) {
	err = dal.DB.QueryRow(`SELECT COUNT(*) FROM transactions WHERE user_id = ? AND dir = ? AND DATE(created_at) = DATE(NOW());`, userID, types.TxTypeWithdrawal).Scan(&count)
	return
}
