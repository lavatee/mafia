package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Superpower struct {
	Name string `json:"name"`
}

func (e *Endpoint) NewSuperpower(c *gin.Context) {
	userId, err := e.GetUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var power Superpower
	err = c.BindJSON(&power)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	powerId, err := e.service.NewSuperpower(userId, power.Name)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"id": powerId,
	})
}
