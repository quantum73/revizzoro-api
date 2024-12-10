package startup

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quantum73/revizzoro-api/config"
	"github.com/quantum73/revizzoro-api/database"
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

	dbConfig := database.DbConfig{
		User:               env.DBUser,
		Password:           env.DBPassword,
		Host:               env.DBHost,
		Port:               env.DBPort,
		DbName:             env.DBName,
		SSLMode:            env.DBSSLMode,
		MaxOpenConnections: env.DbMaxOpenConnections,
		MaxIdleConnections: env.DbMaxIdleConnections,
	}

	db := database.NewDatabase(ctx, dbConfig)
	db.Connect()

	//mainRouter := mux.NewRouter()
	//mainRouter.StrictSlash(true)
	//// Dishes package router
	//dishesRouter := mainRouter.PathPrefix("/dishes").Subrouter()
	//dishesRouter.HandleFunc("/{id:[0-9]+}", dishes.DetailByIdHandler)
	//dishesRouter.HandleFunc("", dishes.ListHandler)
	//// Restaurants package router
	//restaurantsRouter := mainRouter.PathPrefix("/restaurants").Subrouter()
	//restaurantsRouter.HandleFunc("/{id:[0-9]+}", restaurants.DetailByIdHandler)
	//restaurantsRouter.HandleFunc("", restaurants.ListHandler)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

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

	// Disconnect db
	db.Disconnect()
	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error:", err)
	}

	select {
	case <-ctx.Done():
		log.Infof("Gracefully shutting down after %d seconds", env.ServerGracefulTimeout)
	}
}
