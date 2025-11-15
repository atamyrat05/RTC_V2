package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"time"

	"server/internal/configs"
	handler "server/internal/delivery/manager"
	"server/internal/delivery/ws"
	"server/internal/helper"
	spsql "server/pkg/db"
	slog "server/pkg/logging"
	"syscall"

	"github.com/rs/cors"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := configs.GetConfig()
	logger := slog.GetLogger(cfg.Log.Path, cfg.Log.Filename)

	psqlClient, err := spsql.NewClient(ctx,
		spsql.Options{
			Host:          cfg.Storage.Psql.Host,
			Port:          cfg.Storage.Psql.Port,
			Database:      cfg.Storage.Psql.Database,
			Username:      cfg.Storage.Psql.Username,
			Password:      cfg.Storage.Psql.Password,
			PgPoolMaxConn: cfg.Storage.Psql.PgPoolMaxConn,
		})

	if err != nil {
		logger.Errorf("psql client does not connect: %v", err)
		panic("psql client does not connect")
	}
	defer psqlClient.Close()
	logger.Info("psql client connected")

	hub := ws.NewHub()
	router := handler.Manager(logger, psqlClient, cfg, hub)

	err = helper.GetAllRooms(ctx, psqlClient, hub)
	if err != nil {
		logger.Errorln("Error with Get All Rooms", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		AllowedHeaders: []string{
			"*",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(router)

	srv := &http.Server{
		Addr:              cfg.Listen.Port,
		Handler:           handler,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	go hub.Run()
	go func() {
		logger.Info("Starting the server on port", cfg.Listen.Port)
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("ListenAndServe: %v", err)
		}

	}()

	<-ctx.Done()

	logger.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("final")
}
