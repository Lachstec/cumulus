package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/services"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func main() {

	// initialize the database
	db := db_init()

	// initialize the services
	server_service := services.NewServerService(db)
	user_service := services.NewUserService(db)

	// initialize the router
	router := gin.Default()

	// CRUD users
	router.GET("/users", func(c *gin.Context) {
		users, err := user_service.ReadAllUsers()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, users)
	})

	router.POST("/users", func(c *gin.Context) {
		var user types.User
		err := c.BindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		userid, err := user_service.CreateUser(user)
		if err != nil {
			c.AbortWithError(http.StatusConflict, err)
		}
		c.JSON(http.StatusOK, userid)
	})

	router.GET("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		user, err := user_service.ReadUserByUserID(userid)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		c.JSON(http.StatusOK, user)
	})
	
	router.PATCH("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		var user types.User
		err = c.BindJSON(&user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		user, err = user_service.UpdateUser(userid, user)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		c.JSON((http.StatusOK), user)
	})
	
	router.DELETE("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		err = user_service.DeleteUserByUserID(userid)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		c.Status(http.StatusNoContent)
	})

	// users/:userid/servers
	router.GET("/users/:userid/servers", genericEndpoint)

	// self -> // return user by checking bearer token //TODO
	router.GET("/self", genericEndpoint)

	// CRUD servers
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

	//TODO
	router.POST("/servers/:serverid", genericEndpoint)
	router.PUT("/servers/:serverid", genericEndpoint)

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

	router.DELETE("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		err = server_service.DeleteServerByServerID(serverid)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
		}
		c.Status(http.StatusNoContent)
	})

	// servers/:serverid/health
	router.GET("/servers/:serverid/health", genericEndpoint)

	// teapot
	router.GET("/teapot", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	// run webserver
	router.Run("localhost:10000")
}
