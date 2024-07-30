package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, status int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(status, ErrorResponse{Message: message})
}
