package main

import (
	"fmt"
	"github.com/Lachstec/mc-hosting/internal/openstack"
	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
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
	db := dbInit()
	cfg := config.LoadConfig()
	openstack, err := openstack.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// initialize the services
	server_service := services.NewServerService(db)
	user_service := services.NewUserService(db)
	minecraft_provisioner_service := services.NewMinecraftProvisioner(db, openstack, cfg.CryptoConfig.EncryptionKey)

	// initialize the router
	router := gin.Default()

	// CRUD users
	router.GET("/users", func(c *gin.Context) {
		users, err := user_service.ReadAllUsers()
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
		userid, err := user_service.CreateUser(user)
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
		users, err := user_service.ReadUserByUserID(userid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if len(users) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		}
		user := users[0]
		c.JSON(http.StatusOK, user)
	})

	router.PATCH("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		users, err := user_service.ReadUserByUserID(userid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}

		if len(users) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		}

		user := users[0]
		user.ID = userid
		err = c.BindJSON(&user)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		updated, err := user_service.UpdateUser(user)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, updated)
	})

	router.DELETE("/users/:userid", func(c *gin.Context) {
		userid, err := urlParamToInt64(c.Param("userid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		users, err := user_service.ReadUserByUserID(userid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if len(users) == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		user := users[0]
		err = user_service.DeleteUser(user)
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
		servers, err := server_service.ReadAllServers()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, servers)
	})

	//TODO server auch im Openstack erstellen
	router.POST("/servers", func(c *gin.Context) {
		var server *types.Server
		err := c.BindJSON(&server)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		serverid, err := server_service.CreateServer(server)
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
		servers, err := server_service.ReadServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if len(servers) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		}
		server := servers[0]
		c.JSON(http.StatusOK, server)
	})

	//TODO nur gameserver starten, nicht erstellen
	router.POST("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		servers, err := server_service.ReadServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if len(servers) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		}
		server := servers[0]
		if server.Status != types.Stopped {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Server already running/restarting")
		}
		server, err = minecraft_provisioner_service.NewGameServer(c, server)
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
		server.ID = serverid
		err = c.BindJSON(&server)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		server, err = server_service.UpdateServer(server)
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
		servers, err := server_service.ReadServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if len(servers) == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		server := servers[0]
		err = server_service.DeleteServer(server)
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
