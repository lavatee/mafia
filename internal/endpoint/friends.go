package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) GetFriends(c *gin.Context) {
	id, err := e.GetUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	friends, err := e.service.GetFriends(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"friends": friends,
	})
}

func (e *Endpoint) AddFriend(c *gin.Context) {
	userId, err := e.GetUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	friendId, err := strconv.Atoi(c.Param("friend_id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = e.service.AddFriend(userId, friendId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (e *Endpoint) DeleteFriend(c *gin.Context) {

}
