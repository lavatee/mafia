package endpoint

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const tokenKey = "qeq0efquj"

func (e *Endpoint) Middleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	sliceOfHeaders := strings.Split(header, " ")
	if len(sliceOfHeaders) != 2 || sliceOfHeaders[0] != "Bearer" {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid header")
		return
	}
	parsedToken, err := jwt.ParseWithClaims(sliceOfHeaders[1], jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token method")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		c.Set("user_id", claims["id"])
	}
}

func (e *Endpoint) GetUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("invalid getting")
	}
	floatId, ok := id.(float64)
	if !ok {
		return 0, errors.New("invalid type of id")
	}
	return int(floatId), nil
}
