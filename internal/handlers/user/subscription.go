package user

import (
	"fanclub/internal/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) GetTiers(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")

	tiers_data, err := h.dal.Tiers.GetTiersByUsername(username)
	if err != nil {
		h.SendSuccess(c, "failed to get")
		fmt.Printf("filed to get: %v", err)
		return
	}
	h.SendSuccess(c, tiers_data)
}

func (h *HandlerImpl) GetOneTier(c *gin.Context) {
	tierID := c.Request.URL.Query().Get("tierID")
	if tierID == "" {
		h.SendSuccess(c, nil)
		return
	}
	tier_Data, err := h.dal.Tiers.GetOneTierByID(tierID)
	if err != nil {
		fmt.Printf("Error in Get One Tier:%v\n", err)
		h.SendSuccess(c, nil)
		return
	}
	h.SendSuccess(c, tier_Data)
}

func (h *HandlerImpl) CreateTier(c *gin.Context) {
	var req TierBindingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendSuccess(c, "failed to bind")
		fmt.Printf("Error binding: %v\n", err)
		return
	}
	var newTierModel models.TierModel
	newTierModel.Username = req.UserName
	newTierModel.Title = req.TierName
	newTierModel.Color = req.TierColor
	newTierModel.Benefit = req.Benefits
	newTierModel.Child = req.TierChild
	newTierModel.Baseprice = req.BasePrice
	newTierModel.Month2 = req.MonthTwo
	newTierModel.Month3 = req.MonthThree
	newTierModel.Month6 = req.MonthSix
	newTierModel.Active = "1"
	err := h.dal.Tiers.CreateTier(newTierModel)
	if err != nil {
		h.SendSuccess(c, "failed")
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) UpdateTier(c *gin.Context) {
	var req TierBindingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendSuccess(c, "failed to bind")
		fmt.Printf("Error binding: %v\n", err)
		return
	}
	var newTierModel models.TierModel
	newTierModel.Username = req.UserName
	newTierModel.Title = req.TierName
	newTierModel.Color = req.TierColor
	newTierModel.Benefit = req.Benefits
	newTierModel.Child = req.TierChild
	newTierModel.Baseprice = req.BasePrice
	newTierModel.Month2 = req.MonthTwo
	newTierModel.Month3 = req.MonthThree
	newTierModel.Month6 = req.MonthSix
	newTierModel.Active = "1"
	newTierModel.ID = req.TierID
	err := h.dal.Tiers.UpdateTier(newTierModel)
	if err != nil {
		h.SendSuccess(c, "failed")
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) ToggleTier(c *gin.Context) {
	tierID := c.Request.URL.Query().Get("id")
	err := h.dal.Tiers.ToggleTier(tierID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) SubscribeTier(c *gin.Context) {
	tierid := c.Request.URL.Query().Get("tierid")
	duration := c.Request.URL.Query().Get("duration")
	userid := c.Request.URL.Query().Get("userid")
	base_price, err := h.dal.Subscription.GetBasePrice(tierid)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}

	_duration, _ := strconv.ParseFloat(duration, 64)
	pay_amount := base_price * _duration
	_userid, _ := strconv.Atoi(userid)
	wallet_credit, err := h.dal.Wallet.GetWalletBalance(_userid)
	if err != nil {
		fmt.Printf("Error getting wallet balance %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}

	if pay_amount > wallet_credit {
		h.SendSuccess(c, "low balance")
		return
	}
	err = h.dal.Wallet.DeductUserCredit(_userid, pay_amount)
	if err != nil {
		fmt.Printf("Error deducting user credit %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	tierOwner, _ := h.dal.Subscription.GetTierOwner(tierid)

	earningFee := h.troner.GetEarningFee()
	fmt.Printf("---------------------fffffeeeeeeee%v\n -----------------------", earningFee)
	err = h.dal.Wallet.AddUserCredit(tierOwner, pay_amount*(100-earningFee)/100)
	if err != nil {
		fmt.Printf("Error adding user credit %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}

	err = h.dal.Subscription.CreateSubscriber(userid, tierid, duration)
	if err != nil {
		fmt.Printf("Error Creating Subscribers %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	err = h.dal.Userdata.CreatePurchase(userid, "subscribe", tierid, pay_amount)
	if err != nil {
		fmt.Printf("Error Creating Purchase %v\n", err)
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, "success")
}

func (h *HandlerImpl) GetSubscriptions(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")
	filter := c.Request.URL.Query().Get("filter")

	// filter 0:Active 1:Expired 2:ALL
	res, err := h.dal.Subscription.GetTierCards(userid, filter)
	if err != nil {
		h.SendSuccess(c, "failed")
		return
	}
	h.SendSuccess(c, res)
}

func (h *HandlerImpl) GetSubscriptionsCount(c *gin.Context) {
	userid := c.Request.URL.Query().Get("userid")

	res, _ := h.dal.Subscription.GetTierCounts(userid)
	h.SendSuccess(c, res)
}

func (h *HandlerImpl) GetApplicationStatus(c *gin.Context) {
	username := c.Request.URL.Query().Get("username")
	res, err := h.dal.Userdata.GetApplicationStatus(username)
	if err != nil {
		h.SendSuccess(c, 3)
		return
	}
	h.SendSuccess(c, res)
}
