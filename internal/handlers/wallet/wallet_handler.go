package wallet

import (
	"fanclub/internal/dal"
	"fanclub/internal/handlers"
	"fanclub/internal/interfaces"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Withdraw(*gin.Context)
	GetBalance(*gin.Context)
	CheckDepositTxs(*gin.Context)
	WithdrawHouseWallet(*gin.Context)
	GetTransactions(*gin.Context)
	GetUserBalance(*gin.Context)
}

type HandlerImpl struct {
	handlers.BaseHandler
	dal    dal.AppDAL
	troner interfaces.Troner
}

var _ Handler = (*HandlerImpl)(nil)

func NewHandler(dal dal.AppDAL, troner interfaces.Troner) Handler {
	return &HandlerImpl{
		dal:    dal,
		troner: troner,
	}
}
