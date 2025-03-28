package main

import (
	"context"
	"github.com/Lachstec/mc-hosting/internal/api"
	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/logging"
	"github.com/Lachstec/mc-hosting/internal/openstack"
	"github.com/Lachstec/mc-hosting/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func dbInit(cfg config.DbConfig, logger *logrus.Logger) *sqlx.DB {
	s, err := sqlx.Open("pgx", cfg.ConnectionURI())
	ctx := context.Background()
	if err != nil {
		logger.WithContext(ctx).Fatal("failed to connect to backend database")
	}
	mig := db.NewMigrator(s)

	err = mig.Migrate("./migrations")
	if err != nil {
		logger.WithContext(ctx).Fatal("failed to create database schema")
	}

	logger.WithContext(ctx).Info("database connected and initialized")
	return s
}

func main() {
	cfg := config.LoadConfig()
	l := logging.Get(*cfg)

	database := dbInit(cfg.Db, l)

	openstack, err := openstack.NewClient(cfg)
	if err != nil {
		l.Fatal("failed to connect to openstack")
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
		UserService:       *userService,
		ServerService:     *serverService,
		Logger:            l,
		FloatingIPService: *floatingIpService,
		Provisioner:       *minecraftProvisionerService,
	}

	api.RegisterRoutes(router, handler)
	_ = router.Run("0.0.0.0:10000")

}
