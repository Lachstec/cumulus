package main

import (
	Data "data"
	Services "services"
	
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
						c.JSON(http.StatusOK, Services.ReadAllServers())
					case method == "POST":
						var server Data.Server
						server.ID = Services.ReadNumOfServers() + 1
						if err:=c.BindJSON(&server);err!=nil{
							c.AbortWithError(http.StatusBadRequest,err)
							return
						}
						Services.CreateServer(server)
						c.Status(http.StatusCreated)
				}
			} else {
				serverid := urlParamToInteger(c, c.Param("serverid")) - 1
				switch
				{
					case method == "GET":
						if serverid <= Services.ReadNumOfServers() {
							c.JSON(http.StatusOK, Services.ReadServerByServerID(serverid))
						} else {
							c.AbortWithStatus(http.StatusNotFound)
						}
					case method == "POST":
						c.JSON(http.StatusOK, "Starting servers is not implemented yet")
					case method == "PUT":
						c.JSON(http.StatusOK, "Shutting down servers is not implemented yet")
					case method == "PATCH":
						var server Data.Server
						if serverid <= Services.ReadNumOfServers() {
							if err:=c.BindJSON(&server);err!=nil{
								c.AbortWithError(http.StatusBadRequest,err)
								return
							}
							//hier Service.UpdateServer aufrufen und den vorhandenen Server mit dem neuen Server updaten
							c.JSON(http.StatusOK, server)
						} else {
							c.AbortWithStatus(http.StatusBadRequest)
						}
					case method == "DELETE":
						if serverid <= Services.ReadNumOfServers() {
							Services.DeleteServerByServerID(serverid)
						} else {
							c.AbortWithStatus(http.StatusBadRequest)
						}
						c.Status(http.StatusNoContent)
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
	router.POST("/servers/:serverid", genericHandler)
	router.PUT("/servers/:serverid", genericHandler)
	router.PATCH("/servers/:serverid", genericHandler)
	router.DELETE("/servers/:serverid", genericHandler)

	router.GET("/servers/:serverid/health", genericEndpoint)

	router.GET("/teapot", func(c *gin.Context) { c.Status(http.StatusTeapot) })

    router.Run("localhost:10000")
}