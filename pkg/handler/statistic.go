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
	topDebtor, err := h.services.Debt.SelectTop3Debtors(myId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can not get top debtors")
		return
	}
	topLender, err := h.services.Debt.SelectTop3Lenders(myId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can not get top lenders")
		return
	}
	if topDebtor == nil {
		commonStatistic.TopDebtor = 0
	} else {
		commonStatistic.TopDebtor = topDebtor[0]
	}
	if topLender == nil {
		commonStatistic.TopLender = 0
	} else {
		commonStatistic.TopLender = topLender[0]
	}

	FriendsDebts, err := h.services.Debt.GetAllDebts(&myId, nil, strconv.Itoa(data.DebtStatusActive), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can not get friends debts")
		return
	}

	commonStatistic.FriendsActiveDebtsNumber = len(FriendsDebts)

	for _, debt := range FriendsDebts {
		commonStatistic.FriendsActiveDebtsAmount += debt.Amount
	}

	myDebts, err := h.services.Debt.GetAllDebts(nil, &myId, strconv.Itoa(data.DebtStatusActive), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can not det my debt")
		return
	}

	commonStatistic.MyActiveDebtsNumber = len(myDebts)

	for _, debt := range myDebts {
		commonStatistic.MyActiveDebtsAmount += debt.Amount
	}

	c.JSON(http.StatusOK, commonStatistic)
}

func (h *Handler) premiumStatistic(c *gin.Context) {
	id, _ := c.Get(userCtx)
	myId, _ := id.(int)

	requester, err := h.services.Authorization.GetUserById(myId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can not get user by id")
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
		newErrorResponse(c, http.StatusInternalServerError, "can not get friends debts")
		return
	}

	premiumStatistic.FriendsDebtsNumber = len(FriendsDebts)

	for _, debt := range FriendsDebts {
		premiumStatistic.FriendsDebtsAmount += debt.Amount
	}

	myDebts, err := h.services.Debt.GetAllDebts(nil, &myId, strconv.Itoa(data.DebtStatusActive), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can not get my debt")
		return
	}

	premiumStatistic.MyDebtsNumber = len(myDebts)

	for _, debt := range myDebts {
		premiumStatistic.MyDebtsAmount += debt.Amount
	}

	c.JSON(http.StatusOK, premiumStatistic)
}
