package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// todo add validation check if exist such user, cant add yourself and if you are already friends
func (h *Handler) addFriend(c *gin.Context) {
	friendId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	id, _ := c.Get(userCtx)
	intId, _ := id.(int)

	err = h.services.Friends.AddFriend(intId, friendId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}
