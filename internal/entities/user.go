package entities

import (
	"fanclub/internal/models"
	"fanclub/internal/types"
	"time"
)

type UserInfo struct {
	ID           int                  `json:"id"`
	Email        string               `json:"email"`
	Name         string               `json:"name"`
	Avatar       string               `json:"avatar"`
	Active       types.UserStatusType `json:"active"`
	ReferralCode string               `json:"referral_code"`
	JoinDate     time.Time            `json:"joined_at"`
}

func (ui *UserInfo) FromModel(data models.User) {
	ui.ID = data.ID
	ui.Email = data.Email
	ui.Name = data.Name
	ui.Avatar = data.Avatar
	ui.Active = data.Active
	// ui.ReferralCode = data.ReferralCode
	ui.JoinDate = data.CreatedAt
}
