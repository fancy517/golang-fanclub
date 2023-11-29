package user

import (
	"fanclub/internal/entities"
	"fanclub/internal/types"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) GetAllUsers(c *gin.Context) {
	var filter types.UserFilter
	if err := c.Bind(&filter); err != nil {
		h.SendBadInput(c)
		return
	}
	filter.SetDefault()

	users, total, err := h.dal.User.GetAll(filter)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to query users, query=%v; %w", filter, err))
		return
	}

	infos := make([]entities.UserInfo, len(users))
	for i, user := range users {
		infos[i].FromModel(user)
	}

	h.SendSuccess(c, gin.H{
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
		"users":     infos,
	})
}
