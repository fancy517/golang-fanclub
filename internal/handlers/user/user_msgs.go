package user

import (
	"fanclub/internal/entities"
	"fanclub/internal/types"
)

type LoginRequest struct {
	Email    string `json:"email" `
	Password string `json:"password"`
	IP       string `json:"ip" `
	Location string `json:"location"`
	// Recaptcha string `json:"recaptcha" binding:"required"`
}

type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}

//	type LoginResponse struct {
//		MeResponse
//	}
type LoginResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

type SignupRequest struct {
	Username string `json:"email" binding:"required"`
	Email    string `json:"signupEmail" binding:"required"`
	Password string `json:"password" binding:"required"`
	// Recaptcha    string `json:"recaptcha" binding:"required"`
	// ReferralCode string `json:"referral_code"`
}

type ReturnUserData struct {
	UserID      int    `json:"userid"`
	Username    string `json:"username"`
	Usertype    string `json:"usertype"`
	Success     string `json:"success"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	OtpEnabled  int    `json:"otp_enabled"`
	OtpSecret   string `json:"otp_secret"`
	Wallet      string `json:"wallet"`
	Avatar      string `json:"avatarUrl"`
	Status      string `json:"status"`
}

type ReturnTAccountData struct {
	UserName     string `json:"user"`
	DisplayName  string `json:"display"`
	Availability string `json:"availability"`
	AvatarUrl    string `json:"avatar_url"`
	BannerUrl    string `json:"banner_url"`
	Location     string `json:"location"`
	Active       int    `json:"active"`
	Likes        int    `json:"likes"`
	Followers    int    `json:"followers"`
	Success      string `json:"success"`
}

type UserActivateRequest struct {
	Code string `json:"code" binding:"required"`
}

type UserActivateResponse struct {
	Success bool `json:"success"`
}

type MeResponse struct {
	User       entities.UserInfo `json:"user"`
	Permission bool              `json:"permission,omitempty"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type BlockAccountRequest struct {
	UserID int `json:"user_id" binding:"required"`
	Block  int `json:"block"`
}

type UpdateProfileRequest struct {
	Name   string `json:"name" binding:"required"`
	Avatar string `json:"avatar"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequest struct {
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ResendActivationCodeRequest struct {
	Token string `json:"token" binding:"required"`
}

type ConfirmActivationCodeRequest struct {
	Token          string `json:"token" binding:"required"`
	ActivationCode string `json:"activation_code" binding:"required"`
}

type ConfirmActivationCodeResponse struct {
	Success string `json:"success"`
	Email   string `json:"email"`
}

type UploadFilesRequest struct {
	Username string `json:"username"`
}

type LikePostBindingsRequest struct {
	PostID   string `json:"postid" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type RemovePostBindingsRequest struct {
	PostID string `json:"postid" binding:"required"`
	UserID string `json:"userid" binding:"required"`
}

type TierBindingsRequest struct {
	TierID     string `json:"tierid" binding:"required"`
	UserName   string `json:"username" binding:"required"`
	TierName   string `json:"tier_name" binding:"required"`
	TierColor  string `json:"tier_color" binding:"required"`
	BasePrice  string `json:"base_price" binding:"required"`
	Benefits   string `json:"benefits" `
	TierChild  string `json:"tier_child" `
	MonthTwo   string `json:"month_two" binding:"required"`
	MonthThree string `json:"month_three" binding:"required"`
	MonthSix   string `json:"month_six" binding:"required"`
	Active     string `json:"active" binding:"required"`
}

type DisplayNameRequest struct {
	UserID      int    `json:"userid" binding:"required"`
	Displayname string `json:"displayname" binding:"required"`
}

type SafetyBindingData struct {
	Username    string          `json:"username" binding:"required"`
	Content     int             `json:"content" `
	Locations   []string        `json:"locations"`
	Message     int             `json:"message"`
	Permissions [][]interface{} `json:"permissions"`
}

type RelationShipData struct {
	Following  int     `json:"following"`
	Subscribed string  `json:"subscribed"`
	Tipped     float64 `json:"tipped"`
	Like       int     `json:"like"`
}

type TipCreatorData struct {
	Userid      int     `json:"userid"`
	Creator     int     `json:"creator"`
	Amount      float64 `json:"amount"`
	Description string  `json:"text"`
}

type TipPostData struct {
	Userid      int     `json:"userid"`
	Postid      int     `json:"postid"`
	Amount      float64 `json:"amount"`
	Description string  `json:"text"`
}

type NotificationSettingsData struct {
	UserID     int `json:"userid"`
	IsPush     int `json:"is_push"`
	IsMessage  int `json:"is_message"`
	IsReply    int `json:"is_reply"`
	IsPostlike int `json:"is_postlike"`
	IsFollower int `json:"is_follower"`
}

/*
 * Implement Validation
 */

func (r UserActivateRequest) Validate() bool {
	if len(r.Code) != types.ActivateCodeLen {
		return false
	}
	return ValidateActivationCode(r.Code)
}

func (r UpdatePasswordRequest) Validate() bool {
	return r.NewPassword != ""
}

func (r UpdateProfileRequest) Validate() bool {
	return r.Name != ""
}
