package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) FollowUser(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	followid := c.Request.URL.Query().Get("followid")

	result, err := h.dal.Userdata.FollowUser(userid, followid)
	if err != nil {
		h.SendSuccess(c, "failed to follow user")
		fmt.Print("failed to follow user")
		return
	}
	if result == "add" {
		h.SendSuccess(c, "add")
		return
	} else if result == "remove" {
		h.SendSuccess(c, "remove")
		return
	}
	h.SendSuccess(c, "failed to follow user2")
}

func (h *HandlerImpl) GetRelationship(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	followid := c.Request.URL.Query().Get("followid")

	res := RelationShipData{
		Following:  0,
		Subscribed: "0",
		Tipped:     0,
		Like:       0,
	}

	following, err := h.dal.Userdata.IsFollowing(userid, followid)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
		h.SendSuccess(c, res)
		return
	}
	res.Following = following

	like, err := h.dal.Userdata.IsLike(userid, followid)
	if err != nil {
		fmt.Printf("error1: %v\n", err)
		h.SendSuccess(c, res)
		return
	}
	res.Like = like

	h.SendSuccess(c, res)
}

func (h *HandlerImpl) GetPurchases(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	res, err := h.dal.Userdata.GetPurchases(userid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "")
		return
	}
	h.SendSuccess(c, res)
}

func (h *HandlerImpl) SetOnline(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	err := h.dal.User.SetOnline(userid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) SetAway(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	err := h.dal.User.SetAway(userid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) SetHidden(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	err := h.dal.User.SetHidden(userid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) CreateList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	name := c.Request.URL.Query().Get("name")
	err := h.dal.User.CreateList(userid, name)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) EditList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	name := c.Request.URL.Query().Get("name")
	listid := c.Request.URL.Query().Get("listid")
	err := h.dal.User.EditList(userid, name, listid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) DeleteList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	listid := c.Request.URL.Query().Get("listid")
	err := h.dal.User.DeleteList(userid, listid)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}
