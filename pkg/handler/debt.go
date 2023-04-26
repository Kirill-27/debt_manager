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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
	id, _ := c.Get(userCtx)

	if id != debt.LenderId {
		newErrorResponse(c, http.StatusUnauthorized, "you are not a lender of this debt")
		return
	}
	id, err := h.services.Debt.CreateDebt(debt)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
	var debtorId *int
	if filterDebtor != "" {
		str, err := strconv.Atoi(filterDebtor)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "can not parse debtor_id to int")
			return
		}
		debtorId = &str
	}
	filterLender := c.Query(makeFilter("lender_id"))
	var lenderId *int
	if filterLender != "" {
		str, err := strconv.Atoi(filterLender)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "can not parse lender_id to int")
			return
		}
		lenderId = &str
	}
	var sorts []string

	sortAmount := c.Query("sort")
	if sortAmount != "" {
		sorts = strings.Split(sortAmount, ",")
	}

	debts, err := h.services.Debt.GetAllDebts(debtorId, lenderId, sorts)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusMethodNotAllowed, "this debt is not in pending active status")
		return
	}

	err = h.services.Debt.UpdateStatus(debtId, data.DebtStatusActive)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	debtor, err := h.services.CurrentDebt.GetAllCurrentDebts(&debt.DebtorId, &debt.LenderId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	lender, err := h.services.CurrentDebt.GetAllCurrentDebts(&debt.LenderId, &debt.DebtorId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if debtor != nil {
		newAmountDebtor := debtor[0].Amount + debt.Amount
		err = h.services.CurrentDebt.UpdateAmount(debtor[0].Id, newAmountDebtor)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		newAmountLender := lender[0].Amount - debt.Amount
		err = h.services.CurrentDebt.UpdateAmount(lender[0].Id, newAmountLender)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newCurrenDebtReverse := data.CurrentDebts{
		DebtorID: debt.LenderId,
		LenderId: debt.DebtorId,
		Amount:   -debt.Amount,
	}
	_, err = h.services.CurrentDebt.CreateCurrentDebt(newCurrenDebtReverse)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusMethodNotAllowed, "this debt is not in active status")
		return
	}

	err = h.services.Debt.UpdateStatus(debtId, data.DebtStatusClosed)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func makeFilter(value string) string {
	return "filter[" + value + "]"
}
