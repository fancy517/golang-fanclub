package token

import (
	"fanclub/internal/handlers"
)

type Handler interface{}

type HandlerImpl struct {
	handlers.BaseHandler
}

var _ Handler = (*HandlerImpl)(nil)

func NewHandler() Handler {
	return &HandlerImpl{}
}
