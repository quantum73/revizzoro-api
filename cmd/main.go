package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/quantum73/revizzoro-api/pkg/dishes"
	"github.com/quantum73/revizzoro-api/pkg/restaurants"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var host string
	var port int
	var wait time.Duration

	flag.StringVar(&host, "host", "0.0.0.0", "Host to listen on")
	flag.IntVar(&port, "port", 8000, "Port to listen on")
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"timeout for gracefully wait for existing connections to finish (default 15s)",
	)
	flag.Parse()

	mainRouter := mux.NewRouter()
	mainRouter.StrictSlash(true)
	// Dishes package router
	dishesRouter := mainRouter.PathPrefix("/dishes").Subrouter()
	dishesRouter.HandleFunc("", dishes.RootHandler)
	// Restaurants package router
	restaurantsRouter := mainRouter.PathPrefix("/restaurants").Subrouter()
	restaurantsRouter.HandleFunc("", restaurants.RootHandler)

	// Setting up server with gracefully shutdown
	serverAddr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("Starting server on %s", serverAddr)
	srv := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mainRouter,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownChannel

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	_ = srv.Shutdown(ctx)
	log.Println("Gracefully shutting down")
	os.Exit(0)
}
