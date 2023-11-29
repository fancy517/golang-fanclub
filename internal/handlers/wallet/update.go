package wallet

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) CheckDepositTxs(c *gin.Context) {
	_userID := c.Request.URL.Query().Get("userid")
	userID, _ := strconv.Atoi(_userID)
	fmt.Printf("userid %v\n", userID)
	wallet, err := h.dal.Wallet.GetOne(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get wallet userid(%d); %w", userID, err))
		return
	}

	txs, err := h.troner.GetDepositTxs(wallet.DepositAddress, wallet.LastDeposit+1000)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get deposit txs, userid: %d, wallet: %v; %w", userID, wallet.DepositAddress, err))
		return
	}

	if len(txs) == 0 {
		h.SendSuccess(c, CheckDepositTxsResponse{wallet.Credit, false})
		return
	}

	var totalAmount float64
	for _, deposit := range txs {
		totalAmount += deposit.Amount
	}

	minDeposit := h.troner.GetMinDepositValue()
	if totalAmount < minDeposit {
		h.SendBadRequest(c, fmt.Sprintf("Your deposit is below minimum value(%v TRX)", minDeposit))
		return
	}
	fmt.Printf("%#v\n", txs)

	txID, err := h.troner.Collect2House(wallet.PrivateKey, totalAmount)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to collect, user:%d, wallet:%v; amount:%v; %w", userID, wallet.DepositAddress, totalAmount, err))
		return
	}

	if err := h.dal.Transaction.AddCollectHistory(userID, txID, wallet.DepositAddress, totalAmount); err != nil {
		h.SendError(c, fmt.Errorf("failed to add collect history, tx:%v, user:%d; %w", txID, userID, err))
		return
	}

	var lastDeposit int64
	for _, deposit := range txs {
		if err := h.dal.Transaction.AddDepositHistory(userID, deposit.TxHash, wallet.DepositAddress, deposit.Amount); err != nil {
			h.SendError(c, fmt.Errorf("failed to add deposit history, deposit: %#v; %w", deposit, err))
			return
		}
		if deposit.Timestamp > lastDeposit {
			lastDeposit = deposit.Timestamp
		}
	}

	if err := h.dal.Wallet.UpdateLastDepositTime(userID, lastDeposit); err != nil {
		h.SendError(c, fmt.Errorf("faild to update last deposit time, userid:%d; %w", userID, err))
		return
	}

	if err := h.dal.Wallet.AddUserCredit(userID, totalAmount); err != nil {
		h.SendError(c, fmt.Errorf("faild to update user credit, userid:%d, amount:%v; %w", userID, totalAmount, err))
		return
	}

	res := CheckDepositTxsResponse{
		Balance:    wallet.Credit + totalAmount,
		NewDeposit: len(txs) > 0,
	}
	h.SendSuccess(c, res)
}
