package main

import (
	Data "data"
	"strconv"

	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
) 

func genericEndpoint(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Not implemented yet") 
}

func getServersByUserID(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("Not implemented / given userID: %s",  c.Param("userID")))
}

func urlParamToInteger(c *gin.Context, param string) int {
	i, err := strconv.Atoi(param)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest,err)
		panic(err)
	}
	return i
}

func genericHandler(c *gin.Context) {
	
	path := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
	method := c.Request.Method
	switch 
	{
		case path[0] == "servers":
			if len(path) == 1 {
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
						c.JSON(http.StatusCreated,"")
				}
			} else {
				switch
				{
					case method == "GET":
						serverid := urlParamConverter(c, c.Param("serverid")) - 1
						if serverid <= len(Data.Servers) {
							c.IndentedJSON(http.StatusOK, Data.Servers[serverid])
						} else {
							c.AbortWithStatus(http.StatusNotFound)
						}
					case method == "POST":
					case method == "PUT":
					case method == "PATCH":
					case method == "DELETE":
				}
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

	router.GET("/servers/:serverid", genericHandler)
	router.POST("/servers/:serverid", genericEndpoint)
	router.PUT("/servers/:serverid", genericEndpoint)
	router.PATCH("/servers/:serverid", genericEndpoint)
	router.DELETE("/servers/:serverid", genericEndpoint)

	router.GET("/servers/:serverid/health", genericEndpoint)

    router.Run("localhost:10000")
}