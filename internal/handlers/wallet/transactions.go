package wallet

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) GetTransactions(c *gin.Context) {
	_userID := c.Request.URL.Query().Get("userid")
	userID, _ := strconv.Atoi(_userID)
	transactions, err := h.dal.Transaction.GetTransactions(userID)
	if err != nil {
		fmt.Printf("error : %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, transactions)
}
