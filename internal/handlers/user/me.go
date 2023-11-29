package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) GetMe(c *gin.Context) {
	userID, isAnonymous := h.GetUserIDFromContext(c)
	if isAnonymous {
		h.SendNonAuthorized(c)
		return
	}

	user, err := h.dal.User.GetByID(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get user, user_id=%d; %w", userID, err))
		return
	}

	// isAdmin := h.dal.Perm.IsAuthorized(userID)

	var res MeResponse
	res.User.FromModel(user)
	res.Permission = false
	h.SendSuccess(c, res)
}
