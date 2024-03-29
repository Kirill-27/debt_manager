package handler

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/helpers"
	"github.com/kirill-27/debt_manager/requests"

	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
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
	err := validation.Validate(input.Email, validation.Required, is.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	input.Password = helpers.GeneratePasswordHash(input.Password)
	user, err := h.services.Authorization.GetUser(&input.Email, nil)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get user by email")
		return
	}
	if user != nil {
		newErrorResponse(c, http.StatusBadRequest, "user with this email address already exists")
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not create user")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
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
	var input requests.SignInInput

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
		"token": token,
		"id":    userId,
	})
}

func (h *Handler) getAllUsers(c *gin.Context) {
	filterMyFriends := c.Query(makeFilter("my_friends"))
	var friendsFor *int
	if filterMyFriends == "true" {
		id, _ := c.Get(userCtx)
		intId, _ := id.(int)
		friendsFor = &intId
	}

	var sorts []string
	sortAmount := c.Query("sort")
	if sortAmount != "" {
		sorts = strings.Split(sortAmount, ",")
	}

	users, err := h.services.Authorization.GetAllUsers(sorts, friendsFor)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get all users")
		return
	}

	id, _ := c.Get(userCtx)
	idValue, _ := id.(int)
	requester, err := h.services.Authorization.GetUserById(idValue)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get user by id")
		return
	}
	if requester == nil {
		newErrorResponse(c, http.StatusBadRequest, "wrong id in auth token")
		return
	}

	for index := range users {
		users[index].Password = ""
	}
	if requester.SubscriptionType == data.SubscriptionTypeFree {
		for index := range users {
			users[index].Rating = 0
			users[index].MarksSum = 0
			users[index].MarksNumber = 0
		}
	}
	if users == nil {
		c.JSON(http.StatusOK, []data.User{})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) updateUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	id, _ := c.Get(userCtx)
	if id != userId {
		newErrorResponse(c, http.StatusMethodNotAllowed, "you can not change info about user with such id")
		return
	}

	user, err := h.services.Authorization.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get user by id")
		return
	}
	if user == nil {
		newErrorResponse(c, http.StatusNotFound, "there is not user with such id")
		return
	}

	var fieldsToUpdate requests.UpdateUser
	if err := c.BindJSON(&fieldsToUpdate); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if fieldsToUpdate.Email != "" && fieldsToUpdate.Email != user.Email {
		err := validation.Validate(fieldsToUpdate.Email, validation.Required, is.Email)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		userEmail, err := h.services.Authorization.GetUser(&fieldsToUpdate.Email, nil)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get user by email")
			return
		}
		if userEmail != nil {
			newErrorResponse(c, http.StatusBadRequest, "user with this email address already exists")
			return
		}

		user.Email = fieldsToUpdate.Email
	}

	if fieldsToUpdate.Photo != "" {
		user.Photo = fieldsToUpdate.Photo
	}
	if fieldsToUpdate.FullName != "" {
		user.FullName = fieldsToUpdate.FullName
	}
	if fieldsToUpdate.Password != "" {
		user.Password = helpers.GeneratePasswordHash(fieldsToUpdate.Password)
	}

	err = h.services.Authorization.UpdateUser(*user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not update user")
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) getUserById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := h.services.Authorization.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get user by id")
		return
	}
	if user == nil {
		newErrorResponse(c, http.StatusNotFound, "user with this id was not found")
		return
	}
	id, _ := c.Get(userCtx)
	user.Password = ""
	if id != user.Id {
		idValue, _ := id.(int)
		requester, err := h.services.Authorization.GetUserById(idValue)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not get user by id")
			return
		}
		if requester == nil {
			newErrorResponse(c, http.StatusBadRequest, "wrong id in auth token")
			return
		}
		if requester.SubscriptionType == data.SubscriptionTypeFree {
			user.Rating = 0
			user.MarksNumber = 0
			user.MarksSum = 0
		}
	}
	c.JSON(http.StatusOK, *user)
}
