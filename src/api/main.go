package main

import (
	Data "data"
	Services "services"
	
	"strconv"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
) 

func genericEndpoint(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Not implemented yet") 
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
				if len(path) == 3 {
					c.JSON(http.StatusOK, "Healthcheck is not implemented yet")
				} else {
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
								Services.UpdateServer(serverid, server)
								c.JSON(http.StatusOK, Services.ReadServerByServerID(serverid))
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
			}
		case path[0] == "users":
			if len(path) == 1 {
				switch  
				{
					case method == "GET":
						c.JSON(http.StatusOK, Services.ReadAllUsers())
					case method == "POST":
						var user Data.User
						user.ID = Services.ReadNumOfUsers() + 1
						if err:=c.BindJSON(&user);err!=nil{
							c.AbortWithError(http.StatusBadRequest,err)
							return
						}
						Services.CreateUser(user)
						c.Status(http.StatusCreated)
				}
			} else {
				userid := urlParamToInteger(c, c.Param("userid")) - 1
				if len(path) == 3 {
					//c.JSON(http.StatusOK, "Healthcheck is not implemented yet")
				} else {
					switch
					{
						// TODO else auslagern und vorher checken
						case method == "GET":
							if userid <= Services.ReadNumOfUsers() {
								c.JSON(http.StatusOK, Services.ReadUserByUserID(userid))
							} else {
								c.AbortWithStatus(http.StatusNotFound)
							}
						case method == "PATCH":
							var user Data.User
							if userid <= Services.ReadNumOfUsers() {
								if err:=c.BindJSON(&user);err!=nil{
									c.AbortWithError(http.StatusBadRequest,err)
									return
								}
								Services.UpdateUser(userid, user)
								c.JSON(http.StatusOK, Services.ReadUserByUserID(userid))
							} else {
								c.AbortWithStatus(http.StatusBadRequest)
							}
						case method == "DELETE":
							if userid <= Services.ReadNumOfUsers() {
								Services.DeleteUserByUserID(userid)
							} else {
								c.AbortWithStatus(http.StatusBadRequest)
							}
							c.Status(http.StatusNoContent)
					}
				}
			}

		default:
			break
	}
	
}

func main() {
	router := gin.Default()

	// users
	router.GET("/users", genericHandler)
	router.POST("/users", genericHandler)

	// users/:userid
	router.GET("/users/:userid", genericHandler)
	router.PATCH("/users/:userid", genericHandler)
	router.DELETE("/users/:userid", genericHandler)

	// users/:userid/servers
	router.GET("/users/:userid/servers", genericEndpoint)

	// self -> // return user by checking bearer token
	router.GET("/self", genericEndpoint)

	// servers
	router.GET("/servers", genericHandler)
	router.POST("/servers", genericHandler)

	// servers/:serverid
	router.GET("/servers/:serverid", genericHandler)
	router.POST("/servers/:serverid", genericHandler)
	router.PUT("/servers/:serverid", genericHandler)
	router.PATCH("/servers/:serverid", genericHandler)
	router.DELETE("/servers/:serverid", genericHandler)

	// servers/:serverid/health
	router.GET("/servers/:serverid/health", genericHandler)

	// teapot
	router.GET("/teapot", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	// run webserver
    router.Run("localhost:10000")
}