package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, h *Handler) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("", h.GetUsers)
		userGroup.POST("", h.CreateUser)
		userGroup.GET("/:userid", h.GetUserById)
		userGroup.PATCH("/:userid", h.UpdateUserById)
		userGroup.DELETE("/:userid", h.DeleteUserById)
		userGroup.GET("/:userid/servers", h.ServersOfUser)
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
