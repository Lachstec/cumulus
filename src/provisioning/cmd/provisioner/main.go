package main

import (
	"fmt"
	"log"

	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/openstack"
	"github.com/Lachstec/mc-hosting/internal/services"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gin-contrib/cors"
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
	database := dbInit()
	cfg := config.LoadConfig()
	openstack, err := openstack.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	serverStore := db.NewServerStore(database)
	userStore := db.NewUserStore(database)
	ipStore := db.NewIPStore(database)

	// initialize the services
	serverService := services.NewServerService(serverStore)
	userService := services.NewUserService(userStore)
	floatingIpService := services.NewFloatingIPService(ipStore)
	minecraftProvisionerService := services.NewMinecraftProvisioner(database, openstack, cfg.CryptoConfig.EncryptionKey)

	router := gin.Default()

	router.Use(cors.Default())

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
		users, err := userService.ReadUserByUserID(userid)
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

		users, err := userService.ReadUserByUserID(userid)
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
		updated, err := userService.UpdateUser(user)
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
		users, err := userService.ReadUserByUserID(userid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if len(users) == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		user := users[0]
		err = userService.DeleteUser(user)
		if err != nil {
			_ = c.AbortWithError(http.StatusGone, err)
		}
		c.Status(http.StatusNoContent)
	})

	router.GET("/users/:userid/servers", func(ctx *gin.Context) {
		userid, err := urlParamToInt64(ctx.Param("userid"))
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
		}
		servers, err := serverService.ReadServerByUserID(userid)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, servers)
	})

	router.GET("/servers", func(c *gin.Context) {
		servers, err := serverService.ReadAllServers()
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

		//TODO: Hier muss das Token von Auth0 verarbeitet werden und der passende User rausgesucht.
		user := types.User{
			ID:    1,
			Sub:   "Samplesub",
			Name:  "Sampleuser",
			Class: types.Admin.Value(),
		}

		srv, err := minecraftProvisionerService.NewGameServer(c, server, &user)

		if err != nil {
			log.Println("error creating new game server:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		// serverid, err := serverService.CreateServer(server)
		if err != nil {
			_ = c.AbortWithError(http.StatusConflict, err)
		}
		c.JSON(http.StatusOK, srv.ID)
	})

	router.GET("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		servers, err := serverService.ReadServerByServerID(serverid)
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
		servers, err := serverService.ReadServerByServerID(serverid)
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
			log.Printf("Error: %v\n", err)
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		servers, err := serverService.ReadServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}

		if len(servers) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		}

		server := servers[0]

		err = c.BindJSON(&server)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		srv, err := serverService.UpdateServer(server)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, srv)
	})

	router.DELETE("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		servers, err := serverService.ReadServerByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		if len(servers) == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		server := servers[0]
		err = serverService.DeleteServer(server)
		if err != nil {
			_ = c.AbortWithError(http.StatusGone, err)
		}
		c.Status(http.StatusNoContent)
	})

	// servers/:serverid/health
	router.GET("/servers/:serverid/health", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
		}
		ip, err := floatingIpService.ReadIpByServerID(serverid)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, ip)
	})

	// teapot
	router.GET("/teapot", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	_ = router.Run("0.0.0.0:10000")
}
