package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) JoinRoom(c *gin.Context) {
	id, err := e.GetUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	roomId, err := e.service.JoinRoom(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"room_id": roomId,
	})
}

func (e *Endpoint) LeaveRoom(c *gin.Context) {
	id, err := e.GetUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	err = e.service.LeaveRoom(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
