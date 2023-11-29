package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) TipCreator(c *gin.Context) {
	var req TipCreatorData
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding")
		return
	}

	credit, err := h.dal.Wallet.GetWalletBalance(req.Userid)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	if credit < req.Amount {
		h.SendSuccess(c, "low_amount")
		return
	}

	err = h.dal.User.AddTipHistory(req.Userid, "user", req.Creator, req.Amount, req.Description)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}

	err = h.dal.Wallet.DeductUserCredit(req.Userid, req.Amount)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}

	err = h.dal.Wallet.AddUserCredit(req.Creator, req.Amount)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) TipPost(c *gin.Context) {
	fmt.Printf("Tip post")
	var req TipPostData
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadInput(c)
		fmt.Printf("Error binding")
		return
	}

	credit, err := h.dal.Wallet.GetWalletBalance(req.Userid)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	if credit < req.Amount {
		h.SendSuccess(c, "low_amount")
		return
	}

	err = h.dal.User.AddTipHistory(req.Userid, "post", req.Postid, req.Amount, req.Description)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}

	err = h.dal.Wallet.DeductUserCredit(req.Userid, req.Amount)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}

	creator, err := h.dal.Postlist.GetUserIDFromPostid(req.Postid)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}

	err = h.dal.Wallet.AddUserCredit(creator, req.Amount)
	if err != nil {
		fmt.Printf("%v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}
