package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// todo add validation user can see debt where he is debtor or lender
func (h *Handler) getAllCurrentDebts(c *gin.Context) {
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

	debts, err := h.services.CurrentDebt.GetAllCurrentDebts(debtorId, lenderId, sorts)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, debts)
}
