package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kirill-27/debt_manager/data"
	"net/http"
)

// todo make validation if there are closed debt for those 2 users
func (h *Handler) createReview(c *gin.Context) {
	var review data.Review
	if err := c.BindJSON(&review); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := c.Get(userCtx)
	idValue, _ := id.(int)
	review.ReviewerId = idValue

	id, err := h.services.Review.CreateReview(review)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// todo add validation user can see debt where he is debtor or lender
func (h *Handler) getAllReviews(c *gin.Context) {

}
