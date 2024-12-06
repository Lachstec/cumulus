package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/services"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func db_init() *sqlx.DB {
	cfg := config.LoadConfig()
	s, err := sqlx.Open("pgx", cfg.Db.ConnectionURI())
	if err != nil {
		panic(err)
	}
	mig := db.NewMigrator(s)

	err = mig.Migrate("../../migrations")
	if err != nil {
		panic(err)
	}

	fmt.Println("typesbase schema has been created")
	return s
}

func genericEndpoint(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Not implemented yet")
}

func urlParamToInt64(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func genericHandler(c *gin.Context) {

	path := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
	method := c.Request.Method
	switch {
	case path[0] == "servers":
		if len(path) == 1 {
			switch {
			case method == "GET":
				c.JSON(http.StatusOK, services.ReadAllServers())
			case method == "POST":
				var server types.Server
				server.ID = services.ReadNumOfServers() + 1
				if err := c.BindJSON(&server); err != nil {
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
				services.CreateServer(server)
				c.Status(http.StatusCreated)
			}
		} else {
			serverid := urlParamToInteger(c, c.Param("serverid")) - 1
			if len(path) == 3 {
				c.JSON(http.StatusOK, "Healthcheck is not implemented yet")
			} else {
				switch {
				case method == "GET":
					if serverid <= services.ReadNumOfServers() {
						c.JSON(http.StatusOK, services.ReadServerByServerID(serverid))
					} else {
						c.AbortWithStatus(http.StatusNotFound)
					}
				case method == "POST":
					c.JSON(http.StatusOK, "Starting servers is not implemented yet")
				case method == "PUT":
					c.JSON(http.StatusOK, "Shutting down servers is not implemented yet")
				case method == "PATCH":
					var server types.Server
					if serverid <= services.ReadNumOfServers() {
						if err := c.BindJSON(&server); err != nil {
							c.AbortWithError(http.StatusBadRequest, err)
							return
						}
						services.UpdateServer(serverid, server)
						c.JSON(http.StatusOK, services.ReadServerByServerID(serverid))
					} else {
						c.AbortWithStatus(http.StatusBadRequest)
					}
				case method == "DELETE":
					if serverid <= services.ReadNumOfServers() {
						services.DeleteServerByServerID(serverid)
					} else {
						c.AbortWithStatus(http.StatusBadRequest)
					}
					c.Status(http.StatusNoContent)
				}
			}
		}
	case path[0] == "users":
		if len(path) == 1 {
			switch {
			case method == "GET":
				c.JSON(http.StatusOK, services.ReadAllUsers())
			case method == "POST":
				var user types.User
				user.ID = services.ReadNumOfUsers() + 1
				if err := c.BindJSON(&user); err != nil {
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
				services.CreateUser(user)
				c.Status(http.StatusCreated)
			}
		} else {
			userid := urlParamToInteger(c, c.Param("userid")) - 1
			if len(path) == 3 {
				//c.JSON(http.StatusOK, "Healthcheck is not implemented yet")
			} else {
				switch {
				// TODO else auslagern und vorher checken
				case method == "GET":
					if userid <= services.ReadNumOfUsers() {
						c.JSON(http.StatusOK, services.ReadUserByUserID(userid))
					} else {
						c.AbortWithStatus(http.StatusNotFound)
					}
				case method == "PATCH":
					var user types.User
					if userid <= services.ReadNumOfUsers() {
						if err := c.BindJSON(&user); err != nil {
							c.AbortWithError(http.StatusBadRequest, err)
							return
						}
						services.UpdateUser(userid, user)
						c.JSON(http.StatusOK, services.ReadUserByUserID(userid))
					} else {
						c.AbortWithStatus(http.StatusBadRequest)
					}
				case method == "DELETE":
					if userid <= services.ReadNumOfUsers() {
						services.DeleteUserByUserID(userid)
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

	// initialize the database
	db := db_init()

	// initialize the services
	server_service := services.NewServerService(db)

	// initialize the router
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
	router.GET("/servers", func(c *gin.Context) {
		servers, err := server_service.ReadAllServers()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, servers)
	})

	router.POST("/servers", func(c *gin.Context) {
		var server types.Server
		err := c.BindJSON(&server)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		serverid, err := server_service.CreateServer(server)
		if err != nil {
			c.AbortWithError(http.StatusConflict, err)
		}
		c.JSON(http.StatusOK, serverid)
	})

	router.GET("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		server, err := server_service.ReadServerByServerID(serverid)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		c.JSON(http.StatusOK, server)
	})

	router.POST("/servers/:serverid", genericHandler)
	router.PUT("/servers/:serverid", genericHandler)

	router.PATCH("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		var server types.Server
		err = c.BindJSON(&server)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		server, err = server_service.UpdateServer(serverid, server)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, server)
	})
	
	router.DELETE("/servers/:serverid", genericHandler)

	// servers/:serverid/health
	router.GET("/servers/:serverid/health", genericHandler)

	// teapot
	router.GET("/teapot", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	// run webserver
	router.Run("localhost:10000")
}
