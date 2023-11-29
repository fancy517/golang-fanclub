package user

import "github.com/gin-gonic/gin"

func (h *HandlerImpl) MuteUser(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	creator := c.Request.URL.Query().Get("creator")
	status, err := h.dal.Userdata.MuteUser(userid, creator)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, status)
}

func (h *HandlerImpl) BlockUser(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	creator := c.Request.URL.Query().Get("creator")
	status, err := h.dal.Userdata.BlockUser(userid, creator)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, status)
}

func (h *HandlerImpl) VipUser(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	creator := c.Request.URL.Query().Get("creator")
	status, err := h.dal.Userdata.VipUser(userid, creator)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, status)
}

func (h *HandlerImpl) CustomUser(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	creator := c.Request.URL.Query().Get("creator")
	listid := c.Request.URL.Query().Get("listid")
	status, err := h.dal.Userdata.CustomUser(userid, creator, listid)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, status)
}

func (h *HandlerImpl) ReportUser(c *gin.Context) {
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) GetMuteList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	list, err := h.dal.Userdata.GetMuteList(userid)
	if err != nil {
		h.SendSuccess(c, nil)
	}
	h.SendSuccess(c, list)
}

func (h *HandlerImpl) GetBlockList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	list, err := h.dal.Userdata.GetBlockList(userid)
	if err != nil {
		h.SendSuccess(c, nil)
	}
	h.SendSuccess(c, list)
}

func (h *HandlerImpl) GetVipList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	list, err := h.dal.Userdata.GetVipList(userid)
	if err != nil {
		h.SendSuccess(c, nil)
	}
	h.SendSuccess(c, list)
}

func (h *HandlerImpl) GetList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	list, err := h.dal.Userdata.GetList(userid)
	if err != nil {
		h.SendSuccess(c, nil)
	}
	h.SendSuccess(c, list)
}

func (h *HandlerImpl) GetCustomList(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	listid := c.Request.URL.Query().Get("listid")
	list, err := h.dal.Userdata.GetCustomList(userid, listid)
	if err != nil {
		h.SendSuccess(c, nil)
	}
	h.SendSuccess(c, list)
}

func (h *HandlerImpl) GetDefaultEntries(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	list, err := h.dal.Userdata.GetDefaultEntries(userid)
	if err != nil {
		h.SendSuccess(c, nil)
	}
	h.SendSuccess(c, list)
}
