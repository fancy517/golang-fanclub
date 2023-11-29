package user

import (
	"fanclub/internal/models"
	"fanclub/internal/types"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *HandlerImpl) UpdatePassword(c *gin.Context) {
	userID, isAnonymous := h.GetUserIDFromContext(c)
	if isAnonymous {
		h.SendNonAuthorized(c)
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		return
	}

	user, err := h.dal.User.GetByID(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("user not found, id=%d; %w", userID, err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		h.SendBadRequest(c, "Invalid old password")
		return
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), BcryptCost)
	if err != nil {
		h.SendError(c, fmt.Errorf("bcryp generate password error; %w", err))
		return
	}

	if err := h.dal.User.UpdatePassword(userID, string(newPassword)); err != nil {
		h.SendError(c, fmt.Errorf("failed to update password, user_id=%d; %w", userID, err))
		return
	}

	h.SendDefaultSuccess(c)
}

func (h *HandlerImpl) BlockAccount(c *gin.Context) {
	var req BlockAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		return
	}

	user, err := h.dal.User.GetByID(req.UserID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get user, user_id=%d; %w", user.ID, err))
		return
	}

	if user.Active == types.UserStatusUnverified {
		h.SendBadRequest(c, "This account is unverified")
		return
	}

	status := types.UserStatusBlocked
	if req.Block == 0 {
		status = types.UserStatusVerified
	}

	if err := h.dal.User.Block(user.ID, status); err != nil {
		h.SendError(c, fmt.Errorf("failed to block user:%d; %w", user.ID, err))
		return
	}

	h.SendDefaultSuccess(c)
}

func (h *HandlerImpl) UpdateProfile(c *gin.Context) {
	userID, isAnonymous := h.GetUserIDFromContext(c)
	if isAnonymous {
		h.SendNonAuthorized(c)
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		h.SendBadInput(c)
		return
	}

	_, err := h.dal.User.GetByID(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to get user on updateprofile, user_id=%d; %w", userID, err))
		return
	}

	if err := h.dal.User.UpdateProfile(userID, models.User{
		Name:   req.Name,
		Avatar: req.Avatar,
	}); err != nil {
		h.SendError(c, fmt.Errorf("failed to update profile, user_id=%d; %w", userID, err))
		return
	}

	h.SendDefaultSuccess(c)
}
