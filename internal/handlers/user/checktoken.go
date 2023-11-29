package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// func (h *HandlerImpl) GetInfoFromToken(c *gin.Context) {
// 	var req TokenRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.SendBadInput(c)
// 		fmt.Printf("Error binding")
// 		return
// 	}
// 	user, err := h.dal.User.GetByToken(req.Token)
// 	if err != nil {
// 		h.SendSuccess(c, nil)
// 		fmt.Printf("Error getting user")
// 		return
// 	}
// 	returndata := ReturnTAccountData{
// 		DisplayName:  user.DisplayName,
// 		UserName:     user.UserName,
// 		AvatarUrl:    user.AvatarUrl,
// 		BannerUrl:    user.BannerUrl,
// 		Location:     user.Location,
// 		Active:       user.Active,
// 		Likes:        user.Likes,
// 		Followers:    user.Followers,
// 		Availability: user.Status,
// 		Success:      "success",
// 	}
// 	h.SendSuccess(c, returndata)
// }

func (h *HandlerImpl) CheckToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding")
		return
	}
	returndata := ReturnUserData{}
	userID, isExpired, err := h.dal.Token.GetUserFromToken(req.Token)
	if err != nil || isExpired {
		returndata.Success = "failed"
		returndata.Usertype = "guest"
		returndata.UserID = 0
		returndata.Username = ""
		h.SendSuccess(c, returndata)
		return
	} else {
		userType, err := h.dal.User.GetByID(userID)
		if err != nil {
			returndata.Success = "failed"
			returndata.Usertype = "guest"
			returndata.UserID = 0
			returndata.Username = ""
			h.SendSuccess(c, returndata)
			return
		} else {
			returndata.Success = "success"
			returndata.UserID = userType.ID
			returndata.Username = userType.Name
			returndata.Email = userType.Email
			returndata.OtpEnabled = userType.OtpEnabled
			returndata.OtpSecret = userType.OtpSecret
			returndata.Avatar = userType.Avatar
			returndata.DisplayName, _ = h.dal.User.GetDisplayNameByID(userType.ID)
			returndata.Wallet, _ = h.dal.User.GetWalletAddr(userType.ID)
			returndata.Status = userType.Status
			if userType.Active == 0 {
				returndata.Usertype = "needactivate"
			} else if userType.Active == 1 {
				returndata.Usertype = "user"
			} else if userType.Active == 2 {
				returndata.Usertype = "creator"
			} else if userType.Active == 3 {
				returndata.Usertype = "admin"
			} else {
				returndata.Usertype = "guest"
			}
			h.SendSuccess(c, returndata)
			return
		}
	}
}
