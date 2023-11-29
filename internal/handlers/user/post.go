package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) LikePost(c *gin.Context) {
	var req LikePostBindingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendSuccess(c, "failed to bind")
		fmt.Printf("Error binding")
		return
	}
	result, err := h.dal.Postlist.LikePost(req.PostID, req.Username)
	if err != nil {
		h.SendSuccess(c, "failed to like post")
		fmt.Print("failed to like post")
		return
	}
	if result == "add" {
		h.SendSuccess(c, "add")
		return
	} else if result == "remove" {
		h.SendSuccess(c, "remove")
		return
	}
	h.SendSuccess(c, "failed to like post2")
}
