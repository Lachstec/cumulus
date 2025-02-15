package main

import (
	"net/http"
	"strconv"

	"github.com/Lachstec/mc-hosting/internal/api"
	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/logging"
	"github.com/Lachstec/mc-hosting/internal/openstack"
	"github.com/Lachstec/mc-hosting/internal/services"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func dbInit(cfg config.DbConfig, logger zerolog.Logger) *sqlx.DB {
	s, err := sqlx.Open("pgx", cfg.ConnectionURI())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to backend database")
	}
	mig := db.NewMigrator(s)

	err = mig.Migrate("./migrations")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create database schema")
	}

	logger.Info().Msg("database connected and initialized")
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
	cfg := config.LoadConfig()
	l := logging.Get(*cfg)

	database := dbInit(cfg.Db, l)

	openstack, err := openstack.NewClient(cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to connect to openstack")
	}

	serverStore := db.NewServerStore(database)
	userStore := db.NewUserStore(database)
	ipStore := db.NewIPStore(database)

	serverService := services.NewServerService(serverStore)
	userService := services.NewUserService(userStore)
	floatingIpService := services.NewFloatingIPService(ipStore)
	minecraftProvisionerService := services.NewMinecraftProvisioner(database, openstack, cfg.CryptoConfig.EncryptionKey)

	router := gin.Default()

	router.Use(services.CORSMiddleware())
	router.Use(logging.LoggingMiddleware(cfg.LoggingConfig))

	handler := &api.Handler{
		UserService:   *userService,
		ServerService: *serverService,
		Logger:        &l,
	}

	api.RegisterRoutes(router, handler)

	router.GET("/servers", func(c *gin.Context) {
		servers, err := serverService.ReadAllServers()
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve servers from database")
			c.String(http.StatusInternalServerError, "failed to retrieve servers")
			return
		}
		c.JSON(http.StatusOK, servers)
	})

	//TODO server auch im Openstack erstellen
	router.POST("/servers", func(c *gin.Context) {
		var server *types.Server
		err := c.BindJSON(&server)
		if err != nil {
			l.Warn().Err(err).Msg("invalid server payload")
			c.String(http.StatusBadRequest, "invalid server payload")
			return
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
			l.Error().Err(err).Int64("user_id", user.ID).Msg("failed to create game server")
			c.String(http.StatusInternalServerError, "failed to provision game server")
		}
		// serverid, err := serverService.CreateServer(server)
		// if err != nil {
		//	 _ = c.AbortWithError(http.StatusConflict, err)
		// }
		c.JSON(http.StatusOK, srv.ID)
	})

	router.GET("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for server")
			c.String(http.StatusBadRequest, "invalid server format")
			return
		}
		servers, err := serverService.ReadServerByServerID(serverid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve server from database")
			c.String(http.StatusInternalServerError, "failed to retrieve server")
			return
		}
		if len(servers) == 0 {
			c.String(http.StatusNotFound, "no server with given id")
			return
		}
		server := servers[0]
		c.JSON(http.StatusOK, server)
	})

	//TODO nur gameserver starten, nicht erstellen
	router.POST("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for server")
			c.String(http.StatusBadRequest, "invalid server format")
			return
		}
		servers, err := serverService.ReadServerByServerID(serverid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve server from database")
			c.String(http.StatusInternalServerError, "failed to retrieve server")
			return
		}
		if len(servers) == 0 {
			c.String(http.StatusNotFound, "no server with given id")
			return
		}
		server := servers[0]
		if server.Status != types.Stopped {
			c.String(http.StatusBadRequest, "server already running/starting")
			return
		}
		c.JSON(http.StatusOK, server)
	})

	//TODO
	router.PUT("/servers/:serverid", genericEndpoint)

	router.PATCH("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for server")
			c.String(http.StatusBadRequest, "invalid server format")
			return
		}
		servers, err := serverService.ReadServerByServerID(serverid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve server from database")
			c.String(http.StatusInternalServerError, "failed to retrieve server")
			return
		}

		if len(servers) == 0 {
			c.String(http.StatusNotFound, "no server with given id")
			return
		}

		server := servers[0]

		err = c.BindJSON(&server)
		if err != nil {
			l.Warn().Err(err).Msg("invalid server payload")
			c.String(http.StatusUnprocessableEntity, "server payload not valid")
		}
		srv, err := serverService.UpdateServer(server)
		if err != nil {
			l.Warn().Err(err).Msg("failed to update server from database")
			c.String(http.StatusInternalServerError, "failed to update server")
			return
		}
		c.JSON(http.StatusOK, srv)
	})

	router.DELETE("/servers/:serverid", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for server")
			c.String(http.StatusBadRequest, "invalid server format")
			return
		}
		servers, err := serverService.ReadServerByServerID(serverid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve server from database")
			c.String(http.StatusInternalServerError, "failed to retrieve server")
			return
		}
		if len(servers) == 0 {
			c.String(http.StatusNotFound, "no server with given id")
			return
		}
		server := servers[0]
		err = serverService.DeleteServer(server)
		if err != nil {
			l.Warn().Err(err).Msg("failed to delete server from database")
			c.String(http.StatusInternalServerError, "failed to delete server")
			return
		}
		c.Status(http.StatusNoContent)
	})

	// servers/:serverid/health
	router.GET("/servers/:serverid/health", func(c *gin.Context) {
		serverid, err := urlParamToInt64(c.Param("serverid"))
		if err != nil {
			l.Warn().Err(err).Msg("invalid payload for server")
			c.String(http.StatusBadRequest, "invalid server format")
			return
		}
		ip, err := floatingIpService.ReadIpByServerID(serverid)
		if err != nil {
			l.Warn().Err(err).Msg("failed to retrieve floating ip for server from database")
			c.String(http.StatusInternalServerError, "failed to retrieve server ip address")
			return
		}
		c.JSON(http.StatusOK, ip)
	})

	// teapot
	router.GET("/teapot", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	_ = router.Run("0.0.0.0:10000")

}
