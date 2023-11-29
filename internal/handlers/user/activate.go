package user

import (
	"fanclub/internal/models"
	"fanclub/internal/types"
	"fmt"

	tronaccount "fanclub/pkg/go-tron/account"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) Activate(c *gin.Context) {
	var req UserActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		return
	}

	user, err := h.dal.User.GetByActivationCode(req.Code)
	if err != nil {
		h.SendBadInput(c)
		return
	}

	if user.Active != 0 {
		if user.Active == 1 {
			h.SendBadRequest(c, "Your account is alrady activated. you can log in now")
		} else {
			h.SendBadRequest(c, "Your account is blocked. Contact support team!")
		}
		return
	}

	// Check expiration
	// if user.CreatedAt.Add(time.Hour * 24).Before(time.Now()) {
	// 	h.SendBadRequest(c, "This activation code is expired.")
	// 	return
	// }

	if err := h.dal.User.Activate(user.ID); err != nil {
		h.SendError(c, fmt.Errorf("failed to activate user, code=%v; %w", req.Code, err))
		return
	}

	if err := h.createNewWallet(user.ID); err != nil {
		h.SendError(c, fmt.Errorf("failed to insert new wallet, user_id:%d; %w", user.ID, err))
		return
	}

	h.mailer.SendCongrats(user.Email, user.Name)

	h.SendDefaultSuccess(c)
}

func (h *HandlerImpl) ResendActivationCode(c *gin.Context) {
	var req ResendActivationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	var vcode = generateActivationCode(types.ActivateCodeLen)

	userID, isExpired, err := h.dal.Token.GetUserFromToken(req.Token)
	if err != nil || isExpired {
		h.SendSuccess(c, "failed")
		return
	}

	if err := h.dal.User.UpdateActivationCode(userID, vcode); err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) ConfirmActivationCode(c *gin.Context) {
	var req ConfirmActivationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("err binding: %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	userID, isExpired, err := h.dal.Token.GetUserFromToken(req.Token)
	if err != nil || isExpired {
		fmt.Printf("err expired: %v", err)
		h.SendSuccess(c, "failed")
		return
	}
	isActivated, err := h.dal.User.IsActivated(userID)
	if err != nil {
		fmt.Printf("err activated: %v\n", err)
		h.SendSuccess(c, "failed")
		return
	} else if isActivated {
		h.SendSuccess(c, "already activated")
		return
	}
	res, err := h.dal.User.ConfirmActivationCode(userID, req.ActivationCode)
	if err != nil {
		fmt.Printf("err code: %v\n", err)
		h.SendSuccess(c, "failed")
		return
	} else if res == "failed" {
		h.SendSuccess(c, "failed")
		return
	}
	// Send activation code
	if err := h.createNewWallet(userID); err != nil {
		h.SendError(c, fmt.Errorf("failed to insert new wallet, user_id:%d; %w", userID, err))
		return
	}

	err = h.dal.User.AddNotificationSettings(userID)
	if err != nil {
		h.SendError(c, fmt.Errorf("failed to insert new notification settings, user_id:%d; %w", userID, err))
		return
	}

	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) createNewWallet(userID int) error {
	acc := tronaccount.NewLocalAccount()
	wallet := models.Wallet{
		UserID:         userID,
		Credit:         0,
		DepositAddress: acc.Address().ToBase58(),
		PrivateKey:     acc.PrivateKey(),
	}

	if err := h.dal.Wallet.Insert(wallet); err != nil {
		return err
	}
	return nil
}
