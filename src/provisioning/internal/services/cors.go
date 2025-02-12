package services

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			ctx.Header("Access-Control-Allow-Origin", origin)
		}

		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH, HEAD")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
