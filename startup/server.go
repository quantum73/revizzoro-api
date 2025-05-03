package startup

import (
	"context"
	"errors"
	"fmt"
	"github.com/quantum73/revizzoro-api/config"
	base_handlers "github.com/quantum73/revizzoro-api/handlers/base"
	dishes_handlers "github.com/quantum73/revizzoro-api/handlers/dishes"
	pg "github.com/quantum73/revizzoro-api/internal/postgres"
	dishes_services "github.com/quantum73/revizzoro-api/services/dishes"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupLogger(mode config.ServerModeEnum) {
	if mode == config.RELEASE {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			PrettyPrint:     false,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	}
}

func StartServer() {
	ctx := context.Background()
	env := config.NewEnv(".env", true)

	// Settings up logging
	setupLogger(env.ServerMode)
	// Setting up Postgres
	dbConfig := pg.DbConfig{
		Type:               env.DBHType,
		User:               env.DBUser,
		Password:           env.DBPassword,
		Host:               env.DBHost,
		Port:               env.DBPort,
		DbName:             env.DBName,
		SSLMode:            env.DBSSLMode,
		MaxOpenConnections: env.DbMaxOpenConnections,
		MaxIdleConnections: env.DbMaxIdleConnections,
		QueryTimeout:       time.Duration(env.DBQueryTimeout) * time.Second,
		MigrationsPath:     env.DbMigrationsPath,
	}
	postgresDB := pg.NewDatabase(dbConfig)
	err := postgresDB.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	dbObj := postgresDB.GetInstance()

	// TODO: Setting up handlers and controllers
	rootController := base_handlers.NewRootController(dbObj)
	dishService := dishes_services.NewDishService(dbObj)
	dishController := dishes_handlers.NewDishController(dbObj, dishService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dishes/{id}/{$}", dishController.DetailById)
	mux.HandleFunc("GET /dishes/{$}", dishController.List)
	mux.HandleFunc("POST /dishes/{$}", dishController.Create)
	mux.HandleFunc("GET /{$}", rootController.Home)

	// Setting up server with gracefully shutdown
	serverAddr := fmt.Sprintf("%s:%d", env.ServerHost, env.ServerPort)
	log.Infof("Starting server on %s\n", serverAddr)
	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * time.Duration(env.ServeWriteTimeout),
		ReadTimeout:  time.Second * time.Duration(env.ServeReadTimeout),
		IdleTimeout:  time.Second * time.Duration(env.ServerIdleTimeout),
		Handler:      mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownChannel

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(env.ServerGracefulTimeout))
	defer cancel()

	if err := postgresDB.Disconnect(); err != nil {
		log.Warnf("Error disconnecting database: %v\n", err)
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Warnf("Error shutting down server: %v\n", err)
	}

	<-ctx.Done()
	log.Infof("Gracefully shutting down after %d seconds\n", env.ServerGracefulTimeout)
}
