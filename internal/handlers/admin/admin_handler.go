package admin

import (
	"fanclub/internal/dal"
	"fanclub/internal/handlers"
	"fanclub/internal/mailer"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetWaitingModels(*gin.Context)
	GetApplicationByUsername(*gin.Context)
	AcceptApplication(*gin.Context)
	RejectApplication(*gin.Context)
}

type HandlerImpl struct {
	handlers.BaseHandler
	dal    dal.AppDAL
	mailer mailer.Mailer
}

var _ Handler = (*HandlerImpl)(nil)

func NewHandler(dal dal.AppDAL, mailer mailer.Mailer) Handler {
	return &HandlerImpl{
		dal:    dal,
		mailer: mailer,
	}
}
