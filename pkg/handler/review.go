package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/requests"
	"net/http"
	"strconv"
	"strings"
)

// todo make validation if there are closed debt for those 2 users
func (h *Handler) createReview(c *gin.Context) {
	var review data.Review
	if err := c.BindJSON(&review); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	reviewerId, _ := c.Get(userCtx)
	idValue, _ := reviewerId.(int)
	user, err := h.services.Authorization.GetUserById(review.LenderId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if user == nil {
		newErrorResponse(c, http.StatusBadRequest, "lender with such id does not exist")
		return
	}

	reviews, err := h.services.Review.GetAllReviews(&idValue, &review.LenderId, nil)
	if reviews != nil {
		newErrorResponse(c, http.StatusBadRequest, "review for this user already exist")
		return
	}
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

// todo add validation user can see reviews is he hes subscription
func (h *Handler) getAllReviews(c *gin.Context) {
	filterReviewer := c.Query(makeFilter("reviewer_id"))
	var reviewerId *int
	if filterReviewer != "" {
		str, err := strconv.Atoi(filterReviewer)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "can not parse reviewer_id to int")
			return
		}
		reviewerId = &str
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

	debts, err := h.services.Review.GetAllReviews(reviewerId, lenderId, sorts)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, debts)
}

func (h *Handler) updateReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	review, err := h.services.Review.GetReviewById(reviewId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if review == nil {
		newErrorResponse(c, http.StatusNotFound, "there is not review with such id")
		return
	}

	id, _ := c.Get(userCtx)
	if id != review.ReviewerId {
		newErrorResponse(c, http.StatusMethodNotAllowed, "you can not change this about review")
		return
	}

	var fieldsToUpdate requests.UpdateReview
	if err := c.BindJSON(&fieldsToUpdate); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	review.Comment = fieldsToUpdate.Comment
	review.Rate = fieldsToUpdate.Rate

	err = h.services.Review.UpdateReview(*review)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
