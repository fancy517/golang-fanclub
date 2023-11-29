package wallet

import (
	"fanclub/internal/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) Withdraw(c *gin.Context) {
	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		return
	}

	userID, isAnonymous := h.GetUserIDFromContext(c)
	if isAnonymous {
		h.SendNonAuthorized(c)
		return
	}

	// Check if user exceed maximum number of withdrawal
	todayWithdrawal, err := h.dal.Transaction.GetTodaysWithdrawlCount(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to calculate today withdrawal count; %w", err))
		return
	}

	if todayWithdrawal >= types.MaximumWithdrawalPerDay {
		h.SendBadRequest(c, "Maximum withdrawal number exceeded")
		return
	}

	wallet, err := h.dal.Wallet.GetOne(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get user wallet; userid = %d; %w", userID, err))
		return
	}

	if wallet.Credit < req.Amount || req.Amount <= 0 {
		h.SendBadRequest(c, "Insufficient funds")
		return
	}

	txHash, err := h.troner.Withdraw(req.Address, req.Amount)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to withdraw; amount=%v, address=%v; %w", req.Amount, req.Address, err))
		return
	}

	// Deduct user credit
	if err := h.dal.Wallet.DeductUserCredit(userID, req.Amount); err != nil {
		h.SendError(c, fmt.Errorf("failed to deduct user credit on withdrawal, user_id=%v, amount=%v; %w", userID, req.Amount, err))
		return
	}

	// Add txhash to transactions table
	if err := h.dal.Transaction.AddWithdrawalHistory(userID, txHash, req.Address, req.Amount); err != nil {
		h.SendError(c, fmt.Errorf("failed to add tx history; txhash=%v; %w", txHash, err))
		return
	}

	res := WithdrawResponse{
		TxHash:  txHash,
		Balance: wallet.Credit - req.Amount,
	}
	h.SendSuccess(c, res)
}

func (h *HandlerImpl) WithdrawHouseWallet(c *gin.Context) {
	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		return
	}

	_, isAnonymous := h.GetUserIDFromContext(c)
	if isAnonymous {
		h.SendNonAuthorized(c)
		return
	}
	adminUserid := -1

	balance, err := h.dal.Wallet.GetHouseWallet()
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get house balance; %w", err))
		return
	}

	if balance < req.Amount || req.Amount <= 0 {
		h.SendBadRequest(c, "Insufficient funds")
		return
	}

	txHash, err := h.troner.Withdraw(req.Address, req.Amount)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to withdraw; amount=%v, address=%v; %w", req.Amount, req.Address, err))
		return
	}

	// Add txhash to transactions table
	if err := h.dal.Transaction.AddWithdrawalHistory(adminUserid, txHash, req.Address, req.Amount); err != nil {
		h.SendError(c, fmt.Errorf("failed to add tx history; txhash=%v; %w", txHash, err))
		return
	}

	res := WithdrawResponse{
		TxHash:  txHash,
		Balance: balance - req.Amount,
	}
	h.SendSuccess(c, res)
}
