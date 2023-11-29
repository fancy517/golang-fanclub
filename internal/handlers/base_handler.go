package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
}

const UserIDKey = "userid" // value should be same with /app/middleware -> keyUserID

// return value: user_id, is_anonymous
func (h *BaseHandler) GetUserIDFromContext(c *gin.Context) (int, bool) {
	id := c.GetInt(UserIDKey)
	if id == 0 {
		return 0, true
	}
	return id, false
}

/*
 * Return methods
 */
func (h *BaseHandler) SendError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "Internal Server Error",
	})
	log.Println(err)
}

func (h *BaseHandler) SendBadInput(c *gin.Context) {
	h.SendBadRequest(c, "Bad Input")
}

func (h *BaseHandler) SendNonAuthorized(c *gin.Context) {
	h.SendBadRequest(c, "You are not authorized!")
}

func (h *BaseHandler) SendBadRequest(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": msg,
	})
}

func (h *BaseHandler) SendSuccess(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, obj)
}

func (h *BaseHandler) SendDefaultSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *BaseHandler) ShowErrorLog(err error) {
	log.Println(err)
}
