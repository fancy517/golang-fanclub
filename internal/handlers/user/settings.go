package user

import (
	"fanclub/internal/types"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

func (h *HandlerImpl) ChangeUsername(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	password := c.Request.URL.Query().Get("password")
	newusername := c.Request.URL.Query().Get("newusername")
	status, err := h.dal.User.ChangeUsername(username, password, newusername)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, status)
}

func (h *HandlerImpl) ChangeDisplayname(c *gin.Context) {
	var req DisplayNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding")
		return
	}
	userid := req.UserID
	displayname := req.Displayname
	err := h.dal.User.ChangeDisplayname(userid, displayname)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}
func (h *HandlerImpl) ChangePassword(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	oldpassword := c.Request.URL.Query().Get("oldpassword")
	newpassword := c.Request.URL.Query().Get("newpassword")
	err := h.dal.User.CheckPassword(username, oldpassword)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "password")
		return
	}
	err = h.dal.User.ChangePassword(username, newpassword)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) ChangeEmail(c *gin.Context) {
}

func (h *HandlerImpl) DeleteAccount(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	err := h.dal.User.DeleteAccount(userid)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) CheckPassword(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	password := c.Request.URL.Query().Get("password")
	err := h.dal.User.CheckPassword(username, password)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	email, _ := h.dal.User.GetEmailByUsername(username)
	vcode1 := generateActivationCode(types.ActivateCodeLen)
	err = h.dal.User.RequestVCodeEmail(username, vcode1)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.mailer.SendSettingCode(email, username, vcode1)
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) CheckVCode1(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	vcode1 := c.Request.URL.Query().Get("vcode")
	newemail := c.Request.URL.Query().Get("newemail")
	status, err := h.dal.User.CheckVCode1(username, newemail, vcode1)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	if status == "incorrect" {
		h.SendSuccess(c, "incorrect")
		return
	}
	vcode2 := generateActivationCode(types.ActivateCodeLen)
	err = h.dal.User.RequestVCodeEmail2(username, vcode2)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.mailer.SendActivationCode(newemail, username, vcode2)
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) CheckVCode2(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	vcode2 := c.Request.URL.Query().Get("vcode")
	newemail := c.Request.URL.Query().Get("email")
	fmt.Printf("%v %v\n", username, vcode2)
	status, err := h.dal.User.CheckVCode2(username, vcode2, newemail)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	if status == "incorrect" {
		h.SendSuccess(c, "incorrect")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) Generate2FA(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Fanclub",
		AccountName: username,
		Period:      30,
		Digits:      6,
	})
	if err != nil {
		h.SendSuccess(c, "")
		return
	}
	qrCode := key.Secret()
	if err != nil {
		h.SendSuccess(c, "")
		return
	}

	h.SendSuccess(c, qrCode)
}

func (h *HandlerImpl) OTPVerify(c *gin.Context) {
	otptoken := c.Request.URL.Query().Get("otptoken")
	qrcode := c.Request.URL.Query().Get("qrcode")
	username := c.Request.URL.Query().Get("username")

	valid := totp.Validate(otptoken, qrcode)
	if !valid {
		fmt.Printf("%v\n", valid)
		h.SendSuccess(c, "failed")
		return
	}
	err := h.dal.User.SetTwoFactorEnabled(username, qrcode)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) OTPValidate(c *gin.Context) {
	otptoken := c.Request.URL.Query().Get("otptoken")
	username := c.Request.URL.Query().Get("username")
	otpsecret, err := h.dal.User.GetOtpSecret(username, otptoken)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	valid := totp.Validate(otptoken, otpsecret)
	if !valid {
		fmt.Printf("%v\n", valid)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) OTPDisable(c *gin.Context) {
	otptoken := c.Request.URL.Query().Get("otptoken")
	username := c.Request.URL.Query().Get("username")
	otpsecret, err := h.dal.User.GetOtpSecret(username, otptoken)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	valid := totp.Validate(otptoken, otpsecret)
	if !valid {
		fmt.Printf("%v\n", valid)
		h.SendSuccess(c, "failed")
		return
	}
	err = h.dal.User.SetTwoFactorDisabled(username)
	if err != nil {
		fmt.Printf("error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) SaveSafety(c *gin.Context) {
	var req SafetyBindingData
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding JSON: %v", err)
		return
	}
	ids, err := h.dal.User.SavePermissions(req.Permissions)
	if err != nil {
		h.SendSuccess(c, "")
	}
	_locations := strings.Join(req.Locations, ",")
	_permsids := strings.Join(ids, ",")
	err = h.dal.User.SaveSafety(req.Username, req.Content, _locations, _permsids, req.Message)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) GetSafety(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	res, err := h.dal.User.GetSafety(username)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, res)
}

func (h *HandlerImpl) GetSessions(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	sessions, err := h.dal.User.GetSessions(username)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, sessions)
}

func (h *HandlerImpl) CloseSession(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	sid := c.Request.URL.Query().Get("sid")
	err := h.dal.User.CloseSession(username, sid)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) ClearSessions(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	err := h.dal.User.ClearSessions(username)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) GetNotificationSettings(c *gin.Context) {
	_userid := c.Request.URL.Query().Get("userid")
	userid, err := strconv.Atoi(_userid)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	notification, err := h.dal.User.GetNotificationSettings(userid)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, notification)
}

func (h *HandlerImpl) UpdateNotificationSettings(c *gin.Context) {
	var req NotificationSettingsData
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding")
		return
	}
	err := h.dal.User.UpdateNotificationSettings(req.UserID, req.IsPush, req.IsMessage, req.IsReply, req.IsPostlike, req.IsFollower)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}
