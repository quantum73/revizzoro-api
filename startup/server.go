package startup

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/api/dishes"
	"github.com/quantum73/revizzoro-api/api/restaurants"
	pg "github.com/quantum73/revizzoro-api/arch/postgres"
	"github.com/quantum73/revizzoro-api/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer() {
	ctx := context.Background()
	env := config.NewEnv(".env", true)

	// Setting up Postgres
	dbConfig := pg.DbConfig{
		User:               env.DBUser,
		Password:           env.DBPassword,
		Host:               env.DBHost,
		Port:               env.DBPort,
		DbName:             env.DBName,
		SSLMode:            env.DBSSLMode,
		MaxOpenConnections: env.DbMaxOpenConnections,
		MaxIdleConnections: env.DbMaxIdleConnections,
	}
	db := pg.NewDatabase(ctx, dbConfig)
	db.Connect()

	// Setting up routers
	gin.SetMode(env.ServerMode)
	router := gin.Default()
	// Restaurants package router
	restaurantsRouter := router.Group("/restaurants")
	restaurantsController := restaurants.NewController(db)
	restaurantsController.MountRoutes(restaurantsRouter)
	// Dishes package router
	dishesRouter := router.Group("/dishes")
	dishesController := dishes.NewController(db)
	dishesController.MountRoutes(dishesRouter)

	// Setting up server with gracefully shutdown
	serverAddr := fmt.Sprintf("%s:%d", env.ServerHost, env.ServerPort)
	log.Infof("Starting server on %s", serverAddr)
	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * time.Duration(env.ServeWriteTimeout),
		ReadTimeout:  time.Second * time.Duration(env.ServeReadTimeout),
		IdleTimeout:  time.Second * time.Duration(env.ServerIdleTimeout),
		Handler:      router.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownChannel

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(env.ServerGracefulTimeout))
	defer cancel()

	db.Disconnect()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error:", err)
	}

	<-ctx.Done()
	log.Infof("Gracefully shutting down after %d seconds", env.ServerGracefulTimeout)
}
