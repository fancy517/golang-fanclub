package debug

import (
	"fanclub/internal/dal"
	"fanclub/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Debug(*gin.Context)
}

type HandlerImpl struct {
	handlers.BaseHandler
	dal dal.AppDAL
}

var _ Handler = (*HandlerImpl)(nil)

func NewHandler(dal dal.AppDAL) Handler {
	return &HandlerImpl{
		dal: dal,
	}
}

func (h *HandlerImpl) Debug(c *gin.Context) {}
