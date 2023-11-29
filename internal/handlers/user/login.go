package user

import (
	"fanclub/internal/models"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
// TokenCookieName = "token"
)

func (h *HandlerImpl) Login(c *gin.Context) {
	_DefaultTokenTTL, _ := strconv.Atoi(os.Getenv("DEFAULT_TOKEN_TTL"))
	DefaultTokenTTL := time.Duration(_DefaultTokenTTL) * time.Hour
	ClientDomain := os.Getenv("CLIENT_DOMAIN")
	TokenCookieName := os.Getenv("TOKEN_COOKIE_NAME")
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding")
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

	var user models.User
	var err error
	containFlag := strings.Contains(req.Email, "@")
	if containFlag {
		user, err = h.dal.User.GetByEmail(req.Email)
	} else {
		user, err = h.dal.User.GetByName(req.Email)
	}

	if err != nil {
		fmt.Printf("Error1: %v\n", err)
		h.SendBadRequest(c, req.Email+" User Not Found")
		return
	}

	// if user.Active == 0 {
	// 	h.SendBadRequest(c, "Your account is not email verified!")
	// 	return
	// }

	// if user.Active != 1 {
	// 	h.SendBadRequest(c, "Your account is blocked! Contact to support team please!")
	// 	return
	// }

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.SendSuccess(c, "Invalid password")
		return
	}

	token, err := h.dal.Token.RefreshUserToken(user.ID, DefaultTokenTTL)
	if err != nil {
		h.SendError(c, fmt.Errorf("refresh user token error; %w", err))
		return
	}

	// Set HttpOnly Cookie
	c.SetCookie(TokenCookieName, token, int(DefaultTokenTTL.Seconds()), "/", ClientDomain, false, true)

	// isAdmin := h.dal.Perm.IsAuthorized(user.ID)
	// res := LoginResponse{}
	// res.User.FromModel(user)
	// res.Permission = isAdmin
	err = h.dal.User.InsertSession(user.Name, req.IP, req.Location)
	if err != nil {
		fmt.Printf("Error inserting session :%v\n", err)
	}
	res := LoginResponse{}
	res.Token = token
	res.Success = true
	h.SendSuccess(c, res)
}

func (h *HandlerImpl) Signout(c *gin.Context) {
	ClientDomain := os.Getenv("CLIENT_DOMAIN")
	TokenCookieName := os.Getenv("TOKEN_COOKIE_NAME")
	userID, isAnonymous := h.GetUserIDFromContext(c)
	if isAnonymous {
		h.SendNonAuthorized(c)
		return
	}

	_, _ = h.dal.Token.RefreshUserToken(userID, 0)

	c.SetCookie(TokenCookieName, "", 0, "/", ClientDomain, false, true)
	h.SendSuccess(c, gin.H{})
}
