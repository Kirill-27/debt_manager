package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/requests"
	"net/http"
	"strconv"
)

func (h *Handler) commonStatistic(c *gin.Context) {
	id, _ := c.Get(userCtx)
	myId, _ := id.(int)

	var commonStatistic requests.CommonStatistic
	commonStatistic.TopFriendsInteraction = []requests.TopFriendsInteraction{}

	FriendsDebts, err := h.services.Debt.GetAllDebts(&myId, nil, strconv.Itoa(data.DebtStatusActive), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	commonStatistic.FriendsDebtsNumber = len(FriendsDebts)

	for _, debt := range FriendsDebts {
		commonStatistic.FriendsDebtsAmount += debt.Amount
	}

	myDebts, err := h.services.Debt.GetAllDebts(nil, &myId, strconv.Itoa(data.DebtStatusActive), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	commonStatistic.MyDebtsNumber = len(myDebts)

	for _, debt := range myDebts {
		commonStatistic.MyDebtsAmount += debt.Amount
	}

	c.JSON(http.StatusOK, commonStatistic)
}

func (h *Handler) premiumStatistic(c *gin.Context) {
	id, _ := c.Get(userCtx)
	myId, _ := id.(int)

	requester, err := h.services.Authorization.GetUserById(myId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if requester == nil {
		newErrorResponse(c, http.StatusBadRequest, "wrong id in auth token")
		return
	}
	if requester.SubscriptionType != data.SubscriptionTypePremium {
		newErrorResponse(c, http.StatusForbidden, "you have no premium subscription")
		return
	}

	var premiumStatistic requests.PremiumStatistic

	FriendsDebts, err := h.services.Debt.GetAllDebts(&myId, nil, strconv.Itoa(data.DebtStatusActive)+","+strconv.Itoa(data.DebtStatusClosed), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	premiumStatistic.FriendsDebtsNumber = len(FriendsDebts)

	for _, debt := range FriendsDebts {
		premiumStatistic.FriendsDebtsAmount += debt.Amount
	}

	myDebts, err := h.services.Debt.GetAllDebts(nil, &myId, strconv.Itoa(data.DebtStatusActive), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	premiumStatistic.MyDebtsNumber = len(myDebts)

	for _, debt := range myDebts {
		premiumStatistic.MyDebtsAmount += debt.Amount
	}

	c.JSON(http.StatusOK, premiumStatistic)
}
