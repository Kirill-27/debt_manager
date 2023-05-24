package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) addFriend(c *gin.Context) {
	friendId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	id, _ := c.Get(userCtx)
	if id == friendId {
		newErrorResponse(c, http.StatusBadRequest, "you can not add yourself as a friend")
		return
	}

	newFriend, err := h.services.Authorization.GetUserById(friendId)
	if newFriend == nil {
		newErrorResponse(c, http.StatusBadRequest, "there are no users with such id")
		return
	}

	intId, _ := id.(int)
	ifExist, err := h.services.Friends.CheckIfFriendExists(intId, friendId)
	if ifExist {
		newErrorResponse(c, http.StatusBadRequest, "you already added this friend")
		return
	}

	err = h.services.Friends.AddFriend(intId, friendId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not add friend")
		return
	}

	err = h.services.Friends.AddFriend(friendId, intId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error on the server. contact support. can not add friend")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
