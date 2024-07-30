package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (e *Endpoint) SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := e.service.SignUp(input.Email, input.Name, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})

}

func (e *Endpoint) SignIn(c *gin.Context) {
	var input SignInInput
	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	access, refresh, err := e.service.SignIn(input.Email, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"access":  access,
		"refresh": refresh,
	})
}

type RefreshInput struct {
	Refresh string `json:"refresh"`
}

func (e *Endpoint) Refresh(c *gin.Context) {
	var input RefreshInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	access, refresh, err := e.service.Refresh(input.Refresh)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"access":  access,
		"refresh": refresh,
	})
}
