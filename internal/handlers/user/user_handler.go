package user

import (
	"fanclub/internal/dal"
	"fanclub/internal/handlers"
	"fanclub/internal/interfaces"
	"fanclub/internal/mailer"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Login(*gin.Context)
	CheckToken(*gin.Context)
	// GetInfoFromToken(*gin.Context)
	Signout(*gin.Context)
	Signup(*gin.Context)
	Activate(*gin.Context)
	GetMe(*gin.Context)
	ResendActivationCode(*gin.Context)
	ConfirmActivationCode(*gin.Context)
	GetAllUsers(*gin.Context)
	UpdatePassword(*gin.Context)
	UpdateProfile(*gin.Context)
	BlockAccount(*gin.Context)
	InitiatePasswordReset(*gin.Context)
	ResetPassword(*gin.Context)
	UploadFiles(*gin.Context)

	// Follower(*gin.Context)
	// Likes(*gin.Context)
	// FollowersCnt(*gin.Context)
	// LikesCnt(*gin.Context)
	// SetSubscriptions(*gin.Context)
	// GetSubscriptions(*gin.Context)
	// ExpiredSubscriptions(*gin.Context)
	// RemoveUploads(*gin.Context)
	// SelectFromVault(*gin.Context)
	// RemoveBookmarkItem(*gin.Context)
	// AddBookmarkItem(*gin.Context)
	// UpdateUserData(*gin.Context)
	// SetStatus(*gin.Context)

	// User Handler
	FollowUser(*gin.Context)
	GetRelationship(*gin.Context)
	GetPurchases(*gin.Context)
	GetSubscriptions(*gin.Context)
	GetSubscriptionsCount(*gin.Context)
	GetApplicationStatus(*gin.Context)
	TipCreator(*gin.Context)
	TipPost(*gin.Context)
	MuteUser(*gin.Context)
	VipUser(*gin.Context)
	BlockUser(*gin.Context)
	ReportUser(*gin.Context)
	CustomUser(*gin.Context)
	GetMuteList(*gin.Context)
	GetBlockList(*gin.Context)
	GetVipList(*gin.Context)
	GetList(*gin.Context)
	GetCustomList(*gin.Context)
	GetDefaultEntries(*gin.Context)
	SetOnline(*gin.Context)
	SetAway(*gin.Context)
	SetHidden(*gin.Context)
	CreateList(*gin.Context)
	EditList(*gin.Context)
	DeleteList(*gin.Context)
	// Creator Handler
	CreatePost(*gin.Context)
	RemovePost(*gin.Context)
	LikePost(*gin.Context)
	GetTiers(*gin.Context)
	GetOneTier(*gin.Context)
	CreateTier(*gin.Context)
	UpdateTier(*gin.Context)
	ToggleTier(*gin.Context)
	SaveProfile(*gin.Context)
	SubscribeTier(*gin.Context)
	SubmitApplication(*gin.Context)
	// Settings Handler
	ChangeUsername(*gin.Context)
	ChangeDisplayname(*gin.Context)
	ChangePassword(*gin.Context)
	ChangeEmail(*gin.Context)
	DeleteAccount(*gin.Context)
	CheckPassword(*gin.Context)
	CheckVCode1(*gin.Context)
	CheckVCode2(*gin.Context)
	Generate2FA(*gin.Context)
	OTPVerify(*gin.Context)
	OTPValidate(*gin.Context)
	OTPDisable(*gin.Context)
	SaveSafety(*gin.Context)
	GetSafety(*gin.Context)
	GetSessions(*gin.Context)
	CloseSession(*gin.Context)
	ClearSessions(*gin.Context)
	GetNotificationSettings(*gin.Context)
	UpdateNotificationSettings(*gin.Context)
	// Guest Handler
	GetProfile(*gin.Context)
	GetPostData(*gin.Context)
	GetOnePost(*gin.Context)
	GetTags(*gin.Context)
	GetMediaCounts(*gin.Context)
}

type HandlerImpl struct {
	handlers.BaseHandler
	dal    dal.AppDAL
	mailer mailer.Mailer
	troner interfaces.Troner
}

var _ Handler = (*HandlerImpl)(nil)

func NewHandler(dal dal.AppDAL, mailer mailer.Mailer, troner interfaces.Troner) Handler {
	return &HandlerImpl{
		dal:    dal,
		mailer: mailer,
		troner: troner,
	}
}
