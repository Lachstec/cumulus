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

	serverGroup := router.Group("/servers")
	{
		serverGroup.GET("", h.GetServers)
		serverGroup.POST("", h.CreateServer)
		serverGroup.GET("/:serverid", h.GetServerById)
		serverGroup.POST("/:serverid", h.StartServerById)
		serverGroup.PUT("/:serverid", h.Put)
		serverGroup.PATCH("/:serverid", h.UpdateServerById)
		serverGroup.DELETE("/:serverid", h.DeleteServerById)
		serverGroup.GET("/:serverid/health", h.GetIpByServerId)
	}

	router.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.GET("/teapot", func(c *gin.Context) {
		c.Status(http.StatusTeapot)
	})
}
