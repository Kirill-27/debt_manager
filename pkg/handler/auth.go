package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kirill-27/debt_manager/data"
	"net/http"
	"strconv"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input data.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == errors.New("wrong email or password").Error() {
			status = http.StatusNotFound
		}
		newErrorResponse(c, status, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user_id": userId,
	})
}

func (h *Handler) getAllUsers(c *gin.Context) {

}

func (h *Handler) updateUser(c *gin.Context) {

}

func (h *Handler) getUserById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := h.services.Authorization.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		newErrorResponse(c, http.StatusNotFound, "user with this id was not found")
		return
	}
	// todo not show rating if it requested by not premium user
	id, _ := c.Get(userCtx)
	if id != user.Id {
		user.Password = ""
		user.Email = ""
	}
	c.JSON(http.StatusOK, *user)
}
