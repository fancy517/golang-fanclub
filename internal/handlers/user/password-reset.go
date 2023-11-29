package user

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *HandlerImpl) InitiatePasswordReset(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		return
	}

	user, err := h.dal.User.GetByEmail(req.Email)
	fmt.Printf("%v\n", user)
	if err != nil {
		h.SendSuccess(c, "not_exist")
		return
	}

	code := generatePasswordResetCode()
	if err := h.dal.User.SetPasswordResetCode(user.ID, code); err != nil {
		h.SendError(c, fmt.Errorf("failed to set password reset code; %w", err))
		return
	}

	h.mailer.SendResetPasswordCode(req.Email, user.Name, code)

	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendSuccess(c, "Binding Failed")
		return
	}

	userID, expiry, err := h.dal.User.GetUserIDFromPasswordResetCode(req.Code)
	if err != nil {
		h.SendSuccess(c, "Invalid code")
		return
	}

	if expiry.Before(time.Now()) {
		h.SendSuccess(c, "Password reset token has expired")
		return
	}

	// change password
	newPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), BcryptCost)
	if err != nil {
		h.SendSuccess(c, fmt.Errorf("bcryp generate password error; %w", err))
		return
	}

	if err := h.dal.User.UpdatePassword(userID, string(newPassword)); err != nil {
		h.SendSuccess(c, fmt.Errorf("failed to update password, user_id=%d; %w", userID, err))
		return
	}

	if err := h.dal.User.DeletePasswordResetCode(req.Code); err != nil {
		h.SendSuccess(c, fmt.Errorf("failed to delete password reset code; %w", err))
		return
	}

	h.SendSuccess(c, "success")
}
