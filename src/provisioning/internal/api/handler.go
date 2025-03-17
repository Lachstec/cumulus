package api

import (
	"github.com/Lachstec/mc-hosting/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	UserService       services.UserService
	ServerService     services.ServerService
	Provisioner       services.MinecraftProvisioner
	FloatingIPService services.FloatingIPService
	Logger            *zerolog.Logger
}

func (h *Handler) respondSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Status: "success",
		Data:   data,
	})
}

func (h *Handler) respondError(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, Response{
		Status: "error",
		Error: &Error{
			Message: message,
			Details: details,
		},
	})
}
