package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
) 

func genericEndpoint(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Not implemented yet")
}

func getServersByUserID(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Not implemented / given userID: %s",  c.Param("userID")))
}

func main() {
	router := gin.Default()
    
	router.GET("/test", genericEndpoint)
	
	router.GET("/user/:userID/servers", getServersByUserID)

	router.GET("/user/:userID", genericEndpoint)
	router.PATCH("/user/:userID", genericEndpoint)
	router.DELETE("/user/:userID", genericEndpoint)

	router.GET("/users", genericEndpoint)
	router.POST("/users", genericEndpoint)

	router.POST("/login", genericEndpoint)

	router.GET("/servers", genericEndpoint)
	router.POST("/servers", genericEndpoint)

	router.GET("/server/:serverID", genericEndpoint)
	router.POST("/server/:serverID", genericEndpoint)
	router.PUT("/server/:serverID", genericEndpoint)
	router.PATCH("/server/:serverID", genericEndpoint)
	router.DELETE("/server/:serverID", genericEndpoint)

	router.GET("/server/:serverID/health", genericEndpoint)

    router.Run("localhost:8080")
}
