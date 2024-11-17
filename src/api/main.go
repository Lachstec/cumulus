package main

import (
	Data "data"

	"fmt"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
) 

func genericEndpoint(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Not implemented yet") 
}

func getServersByUserID(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Not implemented / given userID: %s",  c.Param("userID")))
}

func genericHandler(c *gin.Context) {
	
	path := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
	method := c.Request.Method
	switch 
	{
		case path[0] == "servers":
			if(len(path) == 1) {
				switch  
				{
					case method == "GET":
						c.IndentedJSON(http.StatusOK, Data.Servers)
					case method == "POST":
						var server Data.Server
						server.ID = len(Data.Servers) + 1
						if err:=c.BindJSON(&server);err!=nil{
							c.AbortWithError(http.StatusBadRequest,err)
							return
						}
						Data.Servers = append(Data.Servers, server)
				}

			} else if (path[1] == "test"){
				c.IndentedJSON(http.StatusOK, "Bepis")
			}
		case path[0] == "users":

		default:
			break
	}
	
}

func main() {
	router := gin.Default()
	
	router.GET("/users/:userID/servers", getServersByUserID)

	router.GET("/users/:userID", genericEndpoint)
	router.PATCH("/users/:userID", genericEndpoint)
	router.DELETE("/users/:userID", genericEndpoint)

	router.GET("/users", genericEndpoint)
	router.POST("/users", genericEndpoint)

	router.POST("/login", genericEndpoint)

	router.GET("/servers", genericHandler)
	router.POST("/servers", genericHandler)

	router.GET("/server/:serverID", genericEndpoint)
	router.POST("/server/:serverID", genericEndpoint)
	router.PUT("/server/:serverID", genericEndpoint)
	router.PATCH("/server/:serverID", genericEndpoint)
	router.DELETE("/server/:serverID", genericEndpoint)

	router.GET("/server/:serverID/health", genericEndpoint)

    router.Run("localhost:10000")
}