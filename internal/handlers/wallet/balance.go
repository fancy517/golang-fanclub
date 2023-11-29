package wallet

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) GetBalance(c *gin.Context) {
	userID, isAnonymous := h.GetUserIDFromContext(c)
	if isAnonymous {
		h.SendNonAuthorized(c)
		return
	}

	wallet, err := h.dal.Wallet.GetOne(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get wallet, userid=%v; %w", userID, err))
		return
	}

	res := GetBalanceResponse{
		Balance:        wallet.Credit,
		DepositAddress: wallet.DepositAddress,
	}
	h.SendSuccess(c, res)
}

func (h *HandlerImpl) GetUserBalance(c *gin.Context) {
	_userID := c.Request.URL.Query().Get("userid")
	userID, _ := strconv.Atoi(_userID)
	credit, err := h.dal.Wallet.GetWalletBalance(userID)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, 0)
		return
	}

	h.SendSuccess(c, credit)
}
