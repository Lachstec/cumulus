package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, h *Handler) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("", h.GetUsers)
		userGroup.POST("", h.CreateUser)
		userGroup.GET("/:userid", h.GetUserById)
	}
}
