package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kirill-27/debt_manager/data"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) getDebtById(c *gin.Context) {
	debtId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	debt, err := h.services.Debt.GetDebtById(debtId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get debt by id")
		return
	}
	if debt == nil {
		newErrorResponse(c, http.StatusNotFound, "debt with this id was not found")
		return
	}
	id, _ := c.Get(userCtx)
	if debt.LenderId != id && debt.DebtorId != id {
		newErrorResponse(c, http.StatusMethodNotAllowed, "you cannot get information on debt with this id")
		return
	}
	c.JSON(http.StatusOK, *debt)
}

func (h *Handler) createDebt(c *gin.Context) {
	var debt data.Debt
	if err := c.BindJSON(&debt); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	debtor, err := h.services.Authorization.GetUserById(debt.DebtorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get debtor by id")
		return
	}
	if debtor == nil {
		newErrorResponse(c, http.StatusBadRequest, "no debtor with such id")
		return
	}

	userId, _ := c.Get(userCtx)
	idValue, _ := userId.(int)

	if debt.DebtorId == idValue {
		newErrorResponse(c, http.StatusBadRequest, "you can't be a debtor in debt which you created")
		return
	}

	debt.LenderId = idValue

	id, err := h.services.Debt.CreateDebt(debt)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not create debt")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// todo add validation user can see debt where he is debtor or lender
func (h *Handler) getAllDebts(c *gin.Context) {
	//id, _ := c.Get(userCtx)
	filterDebtor := c.Query(makeFilter("debtor_id"))
	filterLender := c.Query(makeFilter("lender_id"))
	statuses := c.Query(makeFilter("status"))

	var sorts []string

	sortAmount := c.Query("sort")
	if sortAmount != "" {
		sorts = strings.Split(sortAmount, ",")
	}

	debts, err := h.services.Debt.GetAllDebts(filterDebtor, filterLender, statuses, sorts)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get all debts")
		return
	}

	if debts == nil {
		c.JSON(http.StatusOK, []data.Debt{})
		return
	}

	c.JSON(http.StatusOK, debts)
}

func (h *Handler) activateDebt(c *gin.Context) {
	debtId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	debt, err := h.services.Debt.GetDebtById(debtId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get debt by id")
		return
	}
	if debt == nil {
		newErrorResponse(c, http.StatusNotFound, "debt with this id was not found")
		return
	}

	id, _ := c.Get(userCtx)
	if debt.DebtorId != id {
		newErrorResponse(c, http.StatusMethodNotAllowed, "you are not a debtor of this debt")
		return
	}

	if debt.Status != data.DebtStatusPendingActive {
		newErrorResponse(c, http.StatusBadRequest, "this debt is not in pending active status")
		return
	}

	err = h.services.Debt.UpdateStatus(debtId, data.DebtStatusActive)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update debt status")
		return
	}

	debtor, err := h.services.CurrentDebt.GetAllCurrentDebts(&debt.DebtorId, &debt.LenderId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get all current debts")
		return
	}
	lender, err := h.services.CurrentDebt.GetAllCurrentDebts(&debt.LenderId, &debt.DebtorId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get all current debts")
		return
	}

	if debtor != nil {
		newAmountDebtor := debtor[0].Amount + debt.Amount
		err = h.services.CurrentDebt.UpdateAmount(debtor[0].Id, newAmountDebtor)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update amount in current debt")
			return
		}

		newAmountLender := lender[0].Amount - debt.Amount
		err = h.services.CurrentDebt.UpdateAmount(lender[0].Id, newAmountLender)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update amount in current debt")
			return
		}
		c.JSON(http.StatusNoContent, nil)
		return
	}

	newCurrenDebt := data.CurrentDebts{
		DebtorID: debt.DebtorId,
		LenderId: debt.LenderId,
		Amount:   debt.Amount,
	}
	_, err = h.services.CurrentDebt.CreateCurrentDebt(newCurrenDebt)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not create current debt")
		return
	}

	newCurrenDebtReverse := data.CurrentDebts{
		DebtorID: debt.LenderId,
		LenderId: debt.DebtorId,
		Amount:   -debt.Amount,
	}
	_, err = h.services.CurrentDebt.CreateCurrentDebt(newCurrenDebtReverse)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not create current debt")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) closeDebt(c *gin.Context) {
	debtId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	debt, err := h.services.Debt.GetDebtById(debtId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get debt by id")
		return
	}
	if debt == nil {
		newErrorResponse(c, http.StatusNotFound, "debt with this id was not found")
		return
	}

	id, _ := c.Get(userCtx)
	if debt.DebtorId != id {
		newErrorResponse(c, http.StatusMethodNotAllowed, "you are not a debtor of this debt")
		return
	}

	if debt.Status != data.DebtStatusActive {
		newErrorResponse(c, http.StatusBadRequest, "this debt is not in active status")
		return
	}

	err = h.services.Debt.UpdateStatus(debtId, data.DebtStatusClosed)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update debt status")
		return
	}

	debtor, err := h.services.CurrentDebt.GetAllCurrentDebts(&debt.DebtorId, &debt.LenderId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get current debt")
		return
	}
	lender, err := h.services.CurrentDebt.GetAllCurrentDebts(&debt.LenderId, &debt.DebtorId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get current debt")
		return
	}

	newAmountDebtor := debtor[0].Amount - debt.Amount
	err = h.services.CurrentDebt.UpdateAmount(debtor[0].Id, newAmountDebtor)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update current debt amount")
		return
	}

	newAmountLender := lender[0].Amount + debt.Amount
	err = h.services.CurrentDebt.UpdateAmount(lender[0].Id, newAmountLender)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update current debt amount")
		return
	}

	debtsForLenderAndDebtor, err := h.services.Debt.GetAllDebts(
		strconv.Itoa(debt.DebtorId), strconv.Itoa(debt.LenderId), strconv.Itoa(data.DebtStatusActive), nil)

	if newAmountDebtor == 0 && newAmountLender == 0 && debtsForLenderAndDebtor == nil {
		err = h.services.CurrentDebt.DeleteCurrentDebt(debtor[0].Id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not delete current debt")
			return
		}
		err = h.services.CurrentDebt.DeleteCurrentDebt(lender[0].Id)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not delete current debt")
			return
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) deleteDebtById(c *gin.Context) {
	debtId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	debt, err := h.services.Debt.GetDebtById(debtId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get debt by id")
		return
	}
	if debt == nil {
		newErrorResponse(c, http.StatusNotFound, "debt with this id was not found")
		return
	}

	id, _ := c.Get(userCtx)
	if debt.DebtorId != id && debt.LenderId != id {
		newErrorResponse(c, http.StatusMethodNotAllowed, "you are not a debtor or lender of this debt")
		return
	}

	if debt.Status != data.DebtStatusPendingActive {
		newErrorResponse(c, http.StatusMethodNotAllowed, "this debt is not in pending active status")
		return
	}

	err = h.services.Debt.DeleteDebt(debtId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not delete debt")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) closeAllWithDebt(c *gin.Context) {
	lenderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	id, _ := c.Get(userCtx)
	debtorId, _ := id.(int)

	debtorCurrentDebts, err := h.services.CurrentDebt.GetAllCurrentDebts(&debtorId, &lenderId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get current debt")
		return
	}
	if debtorCurrentDebts == nil {
		newErrorResponse(c, http.StatusBadRequest, "there is no current debts between you and this lender")
		return
	}
	if debtorCurrentDebts[0].Amount < 0 {
		newErrorResponse(c, http.StatusBadRequest, "this person doesn't owe you. You owe him")
		return
	}

	lenderCurrentDebts, err := h.services.CurrentDebt.GetAllCurrentDebts(&lenderId, &debtorId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get current debt")
		return
	}
	if lenderCurrentDebts == nil {
		newErrorResponse(c, http.StatusBadRequest, "there is no current debts between you and this lender")
		return
	}

	err = h.services.Debt.CloseAllDebts(debtorId, lenderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not close all debts")
		return
	}
	err = h.services.CurrentDebt.DeleteCurrentDebt(debtorCurrentDebts[0].Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not delete current debt")
		return
	}
	err = h.services.CurrentDebt.DeleteCurrentDebt(lenderCurrentDebts[0].Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not delete current debt")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func makeFilter(value string) string {
	return "filter[" + value + "]"
}
