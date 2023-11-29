package app

import (
	"github.com/gin-gonic/gin"
)

func (app *application) initRouters() {
	r := gin.Default()
	r.Use(app.CORSMiddleware())

	//Fanclub
	r.POST("/upload/files", app.userHandler.UploadFiles)

	r.Static("/public", "./uploads")
	// Guest API
	guestapi := r.Group("/api")
	{
		guestapi.GET("/profile", app.userHandler.GetProfile)
		guestapi.GET("/userpost", app.userHandler.GetPostData)
		guestapi.GET("/getonepost", app.userHandler.GetOnePost)
		// guestapi.GET("/getreply", app.userHandler.GetReplyData)
		guestapi.GET("/gettags", app.userHandler.GetTags)
	}

	authapi := r.Group("/auth")
	authapi.Use(app.authenticate)
	{
		authapi.POST("/login", app.userHandler.Login)
		authapi.POST("/checktoken", app.userHandler.CheckToken)
		authapi.POST("/signup", app.userHandler.Signup)
		authapi.POST("/forgot-password", app.userHandler.InitiatePasswordReset)
		authapi.POST("/reset-password", app.userHandler.ResetPassword)
		authapi.POST("/resend_verification_code", app.userHandler.ResendActivationCode)
		authapi.POST("/confirm_verification_code", app.userHandler.ConfirmActivationCode)
		authapi.POST("/saveprofile", app.userHandler.SaveProfile)
		authapi.GET("/getmediacounts", app.userHandler.GetMediaCounts)
		authapi.GET("/subscribe_tier", app.userHandler.SubscribeTier)
		authapi.POST("/submit_applications", app.userHandler.SubmitApplication)

		//settings
		authapi.GET("/change_username", app.userHandler.ChangeUsername)
		authapi.POST("/change_displayname", app.userHandler.ChangeDisplayname)
		authapi.GET("/change_password", app.userHandler.ChangePassword)
		authapi.GET("/delete_account", app.userHandler.DeleteAccount)
		authapi.GET("/change_email", app.userHandler.ChangeEmail)
		authapi.GET("/check_password", app.userHandler.CheckPassword)
		authapi.GET("/check_vcode1", app.userHandler.CheckVCode1)
		authapi.GET("/check_vcode2", app.userHandler.CheckVCode2)
		authapi.GET("/generate_2fa", app.userHandler.Generate2FA)
		authapi.GET("/otp_verify", app.userHandler.OTPVerify)
		authapi.GET("/otp_validate", app.userHandler.OTPValidate)
		authapi.GET("/otp_disable", app.userHandler.OTPDisable)
		authapi.POST("/save_safety", app.userHandler.SaveSafety)
		authapi.GET("/get_safety", app.userHandler.GetSafety)
		authapi.GET("/get_sessions", app.userHandler.GetSessions)
		authapi.GET("/close_session", app.userHandler.CloseSession)
		authapi.GET("/clear_sessions", app.userHandler.ClearSessions)
		authapi.GET("/get_notification_settings", app.userHandler.GetNotificationSettings)
		authapi.POST("/update_notification_settings", app.userHandler.UpdateNotificationSettings)
		// authapi.GET("/get_safety", app.userHandler.GetSafety)

		userapi := authapi.Group("/user")
		{
			userapi.POST("/likepost", app.userHandler.LikePost)
			userapi.GET("/follow", app.userHandler.FollowUser)
			userapi.GET("/getrelationship", app.userHandler.GetRelationship)
			userapi.GET("/get_purchases", app.userHandler.GetPurchases)
			userapi.GET("/getSubscriptions", app.userHandler.GetSubscriptions)
			userapi.GET("/getSubscriptionsCount", app.userHandler.GetSubscriptionsCount)
			userapi.GET("/get_application_status", app.userHandler.GetApplicationStatus)
			userapi.POST("/tip_creator", app.userHandler.TipCreator)
			userapi.POST("/tip_post", app.userHandler.TipPost)
			userapi.GET("/mute_user", app.userHandler.MuteUser)
			userapi.GET("/block_user", app.userHandler.BlockUser)
			userapi.GET("/vip_user", app.userHandler.VipUser)
			userapi.GET("/custom_user", app.userHandler.CustomUser)
			userapi.GET("/report_user", app.userHandler.ReportUser)
			userapi.GET("/get_mutelist", app.userHandler.GetMuteList)
			userapi.GET("/get_blocklist", app.userHandler.GetBlockList)
			userapi.GET("/get_viplist", app.userHandler.GetVipList)
			userapi.GET("/get_list", app.userHandler.GetList)
			userapi.GET("/get_customlist", app.userHandler.GetCustomList)
			userapi.GET("/get_default_entries", app.userHandler.GetDefaultEntries)
			userapi.GET("/setOnline", app.userHandler.SetOnline)
			userapi.GET("/setAway", app.userHandler.SetAway)
			userapi.GET("/setHidden", app.userHandler.SetHidden)
			userapi.GET("/create_list", app.userHandler.CreateList)
			userapi.GET("/edit_list", app.userHandler.EditList)
			userapi.GET("/delete_list", app.userHandler.DeleteList)

		}
		creatorapi := authapi.Group("/creator")
		{
			creatorapi.POST("/post", app.userHandler.CreatePost)
			creatorapi.POST("/removepost", app.userHandler.RemovePost)
			creatorapi.GET("/tier/gettiers", app.userHandler.GetTiers)
			creatorapi.GET("/tier/getone", app.userHandler.GetOneTier)
			creatorapi.POST("/tier/create", app.userHandler.CreateTier)
			creatorapi.POST("/tier/update", app.userHandler.UpdateTier)
			creatorapi.GET("/tier/toggletier", app.userHandler.ToggleTier)
		}
		walletapi := authapi.Group("/wallet")
		{
			walletapi.GET("/balance", app.walletHandler.GetBalance)
			walletapi.GET("/deposit", app.walletHandler.CheckDepositTxs)
			walletapi.POST("/withdraw", app.walletHandler.Withdraw)
			walletapi.GET("/get_transactions", app.walletHandler.GetTransactions)
			walletapi.GET("/get_balance", app.walletHandler.GetUserBalance)

		}
	}

	adminapi := r.Group("/admin")
	adminapi.Use(app.authenticate)
	{
		adminapi.GET("/applications", app.adminHandler.GetWaitingModels)
		adminapi.GET("/getapplication", app.adminHandler.GetApplicationByUsername)
		adminapi.GET("/acceptapplication", app.adminHandler.AcceptApplication)
		adminapi.GET("/rejectapplication", app.adminHandler.RejectApplication)
	}

	v1api := r.Group("/v1")
	v1api.Use(app.authenticate) // middleware setup
	{
		v1api.POST("/auth/signout", app.userHandler.Signout)

		userapi := v1api.Group("/user")
		{
			userapi.GET("/me", app.userHandler.GetMe)
			userapi.PUT("/update-password", app.userHandler.UpdatePassword)
			userapi.PUT("/update-profile", app.userHandler.UpdateProfile)
		}

	}

	r.SetTrustedProxies(nil)

	app.engine = r
}
