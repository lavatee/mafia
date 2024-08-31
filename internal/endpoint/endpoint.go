package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lavatee/mafia/internal/service"
)

type Endpoint struct {
	service *service.Service
}

func NewEndpoint(svc *service.Service) *Endpoint {
	return &Endpoint{
		service: svc,
	}
}

func (e *Endpoint) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		ctx.Writer.Header().Set("Access-Control-Allow-Creditionals", "true")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
	})
	auth := router.Group("/auth")
	{
		auth.POST("/signup", e.SignUp)
		auth.POST("/signin", e.SignIn)
		auth.POST("/refresh", e.Refresh)
	}
	api := router.Group("/api", e.Middleware)
	{
		api.POST("/friends/:friend_id", e.AddFriend)
		api.GET("/friends", e.GetFriends)
		api.GET("/requests", e.GetRequests)
		api.DELETE("/requests/:sender_id", e.RejectRequest)
		api.POST("/requests/:recipient_id", e.NewRequest)
		api.PUT("/friends/:friend_id", e.DeleteFriend)
		api.POST("/joinroom", e.JoinRoom)
		api.POST("/leaveroom", e.LeaveRoom)
	}
	return router
}
