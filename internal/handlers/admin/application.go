package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) GetWaitingModels(c *gin.Context) {

	users, err := h.dal.User.GetWaitingApplications()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, users)
}

func (h *HandlerImpl) GetApplicationByUsername(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	application, err := h.dal.User.GetApplicationByUsername(username)
	if err != nil {
		h.SendSuccess(c, nil)
		return
	}
	h.SendSuccess(c, application)
}

func (h *HandlerImpl) AcceptApplication(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	err := h.dal.User.AcceptApplication(username)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) RejectApplication(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	err := h.dal.User.RejectApplication(username)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")

}
