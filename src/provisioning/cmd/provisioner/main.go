package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/openstack"
	"github.com/Lachstec/mc-hosting/internal/services"
	"github.com/Lachstec/mc-hosting/internal/types"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func dbInit() *sqlx.DB {
	cfg := config.LoadConfig()
	s, err := sqlx.Open("pgx", cfg.Db.ConnectionURI())
	if err != nil {
		panic(err)
	}
	mig := db.NewMigrator(s)

	err = mig.Migrate("./migrations")
	if err != nil {
		panic(err)
	}

	fmt.Println("typesbase schema has been created")
	return s
}

func cfgInit() (*config.Config, error) {
	key, err := base64.StdEncoding.DecodeString("1YRCJE3rUygZv4zXUhBNUf1sDUIszdT2KAtczVYB85c=")
	if err != nil {
		return nil, err
	}
	cfg := &config.Config{
		Db:    config.DbConfig{},
		Auth0: config.Auth0Config{},
		Openstack: config.OpenStackConfig{
			IdentityEndpoint: "<keystone_url>",
			Username:         "<username>",
			Password:         "<password>",
			Domain:           "<domain>",
			TenantName:       "<tenant_name>",
		},
		CryptoConfig: config.CryptoConfig{
			EncryptionKey: key,
		},
	}
	return cfg, nil
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
	conn := dbInit()

	cfg, err := cfgInit()
	if err != nil {
		panic(err)
	}
	ostack, err := openstack.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// initialize the services
	serverService := services.NewServerService(conn)
	userService := services.NewUserService(conn)
	minecraftProvisionerService := services.NewMinecraftProvisioner(conn, ostack, cfg.CryptoConfig.EncryptionKey)

	// initialize the router
	router := gin.Default()

	// CRUD users
	router.GET("/users", func(c *gin.Context) {
		users, err := userService.ReadAllUsers()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, users)
	})

	router.POST("/users", func(c *gin.Context) {
		var user *types.User
		err := c.BindJSON(&user)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		userid, err := userService.CreateUser(user)
		if err != nil {
			_ = c.AbortWithError(http.StatusConflict, err)
		}
		c.JSON(http.StatusOK, userid)
	})

	router.GET("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		user, err := userService.ReadUserByUserID(userid)
		if err != nil {
			_ = c.AbortWithError(http.StatusNotFound, err)
		}
		c.JSON(http.StatusOK, user)
	})

	router.PATCH("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		var user *types.User
		err = c.BindJSON(&user)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		user, err = userService.UpdateUser(userid, user)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, user)
	})

	router.DELETE("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		err = userService.DeleteUserByUserID(userid)
		if err != nil {
			_ = c.AbortWithError(http.StatusGone, err)
		}
		c.Status(http.StatusNoContent)
	})

	// users/:userid/servers
	router.GET("/users/:userid/servers", genericEndpoint)

	/*router.GET("/self", func(c *gin.Context) {
		token := c.GetHeader("Token")
		user, err := auth_service.ValidateToken(token)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		c.JSON(http.StatusOK, user)
	})*/

	// CRUD servers
	router.GET("/servers", func(c *gin.Context) {
		servers, err := serverService.ReadAllServers()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, servers)
	})

	router.POST("/servers", func(c *gin.Context) {
		var server *types.Server
		err := c.BindJSON(&server)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		serverid, err := serverService.CreateServer(server)
		if err != nil {
			_ = c.AbortWithError(http.StatusConflict, err)
		}
		c.JSON(http.StatusOK, serverid)
	})

	router.GET("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		server, err := serverService.ReadServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusNotFound, err)
		}
		c.JSON(http.StatusOK, server)
	})

	//start game server
	router.POST("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		var server *types.Server
		server, err = serverService.ReadServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusNotFound, err)
		}
		if server.Status != types.Stopped {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Server already running/restarting")
		}
		server, err = minecraftProvisionerService.NewGameServer(c, server)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, server)
	})

	//TODO
	router.PUT("/servers/:serverid", genericEndpoint)

	router.PATCH("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		var server *types.Server
		err = c.BindJSON(&server)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		server, err = serverService.UpdateServer(serverid, server)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, server)
	})

	router.DELETE("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		err = serverService.DeleteServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusGone, err)
		}
		c.Status(http.StatusNoContent)
	})

	// servers/:serverid/health
	router.GET("/servers/:serverid/health", genericEndpoint)

	// teapot
	router.GET("/teapot", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	// run webserver
	_ = router.Run("localhost:10000")
}
