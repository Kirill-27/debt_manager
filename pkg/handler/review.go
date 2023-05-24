package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/requests"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) createReview(c *gin.Context) {
	var review data.Review
	if err := c.BindJSON(&review); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if review.Rate < 1 || review.Rate > 5 {
		newErrorResponse(c, http.StatusBadRequest, "rate value must be between 1 and 5 incl.")
		return
	}

	reviewerId, _ := c.Get(userCtx)
	idValue, _ := reviewerId.(int)

	debts, err := h.services.Debt.GetAllDebts(
		strconv.Itoa(idValue), strconv.Itoa(review.LenderId), strconv.Itoa(data.DebtStatusClosed), nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get closed debts")
		return
	}
	if debts == nil {
		newErrorResponse(c, http.StatusBadRequest, "you don't have any closed deals with lander with this id")
		return
	}

	reviews, err := h.services.Review.GetAllReviews(&idValue, &review.LenderId, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get reviews")
		return
	}
	if reviews != nil {
		newErrorResponse(c, http.StatusBadRequest, "review for this user already exist")
		return
	}
	review.ReviewerId = idValue

	lender, err := h.services.Authorization.GetUserById(review.LenderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get user by id")
		return
	}
	if lender == nil {
		newErrorResponse(c, http.StatusBadRequest, "no such lender with this id")
	}
	lender.MarksNumber += 1
	lender.MarksSum += review.Rate
	lender.Rating = float64(lender.MarksSum) / float64(lender.MarksNumber)

	err = h.services.Authorization.UpdateUser(*lender)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update user")
		return
	}

	id, err := h.services.Review.CreateReview(review)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not create review")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// todo add validation user can see reviews is he has subscription
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

	reviews, err := h.services.Review.GetAllReviews(reviewerId, lenderId, sorts)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get all reviews")
		return
	}
	if reviews == nil {
		c.JSON(http.StatusOK, []data.Review{})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func (h *Handler) updateReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	review, err := h.services.Review.GetReviewById(reviewId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get review by id")
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

	if fieldsToUpdate.Rate < 1 || fieldsToUpdate.Rate > 5 {
		newErrorResponse(c, http.StatusBadRequest, "rate value must be between 1 and 5 incl.")
		return
	}

	lender, err := h.services.Authorization.GetUserById(review.LenderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get lender by id")
		return
	}
	if lender == nil {
		newErrorResponse(c, http.StatusInternalServerError, "no such lender with this id")
	}

	lender.MarksSum += fieldsToUpdate.Rate - review.Rate
	lender.Rating = float64(lender.MarksSum) / float64(lender.MarksNumber)
	err = h.services.Authorization.UpdateUser(*lender)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update user")
		return
	}

	review.Comment = fieldsToUpdate.Comment
	review.Rate = fieldsToUpdate.Rate

	err = h.services.Review.UpdateReview(*review)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update review")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
