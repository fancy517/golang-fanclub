package user

import (
	"fanclub/internal/models"
	"fanclub/internal/types"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 15
)

func (h *HandlerImpl) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadRequest(c, "Bad input")
		return
	}

	// passed, err := recaptcha.VerifyRecaptcha(req.Recaptcha)
	// if err != nil {
	// 	h.SendError(c, fmt.Errorf("failed to verify recaptcha; %w", err))
	// 	return
	// }

	// if !passed {
	// 	h.SendBadRequest(c, "reCaptcha check failed!")
	// 	return
	// }

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), BcryptCost)
	if err != nil {
		h.SendError(c, fmt.Errorf("bcrypt generate password error; %w", err))
		return
	}

	// create new user
	data := models.User{
		Name:           req.Username,
		Email:          req.Email,
		PasswordHash:   string(hash),
		ActivationCode: generateActivationCode(types.ActivateCodeLen),
	}

	err = h.dal.User.InsertGetID(data)
	if err != nil {
		fmt.Printf("error1 :%v\n", err)
		if strings.Contains(err.Error(), "Duplicate entry") {
			h.SendSuccess(c, "exist")
		} else {
			h.SendSuccess(c, "failed")
		}
		return
	}

	h.mailer.SendActivationCode(req.Email, req.Username, data.ActivationCode)

	// Add referral if request has referral code
	// if req.ReferralCode != "" {
	// 	referrer, err := h.dal.User.GetByReferralCode(req.ReferralCode)
	// 	if err == nil {
	// 		if err := h.dal.Referral.Insert(models.Referral{
	// 			RefereeID:  int(userID),
	// 			ReferrerID: referrer.ID,
	// 			Reward:     h.referralBonus,
	// 		}); err != nil {
	// 			h.SendError(c, fmt.Errorf("failed to add referral, user_id=%d; %w", userID, err))
	// 			return
	// 		}

	// 		h.mailer.SendReferral(referrer.Email, referrer.Name, req.Email, req.Name, h.minBetAmountForReferral)
	// 	}
	// }

	h.SendSuccess(c, "success")
}
