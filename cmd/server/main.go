package main

import (
	"boilerplate/internal/container"
	"boilerplate/internal/router"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	// _ "net/http/pprof".

	"github.com/spf13/viper"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx := context.Background()

	logger := slog.Default()

	app := container.NewApplication(ctx, logger)

	mux := router.New()

	mux.HandlerFunc(http.MethodGet, "/healthcheck", app.HealthCheckControler.HandleHeathCheck)
	mux.HandlerFunc(http.MethodGet, "/readiness", app.HealthCheckControler.HandleReadiness)
	mux.HandlerFunc(http.MethodGet, "/cars", app.CarsControler.HandleGetCars)
	mux.HandlerFunc(http.MethodGet, "/cars/:id", app.CarsControler.HandleGetCarByID)
	mux.HandlerFunc(http.MethodPost, "/cars", app.CarsControler.HandleCreateCar)
	mux.HandlerFunc(http.MethodPut, "/cars/:id", app.CarsControler.HandleUpdateCar)

	port := fmt.Sprintf(":%s", viper.GetString("port"))
	srv := &http.Server{
		ReadHeaderTimeout: shutdownTimeout,
		Addr:              port,
		Handler:           mux,
	}

	logger.Info("starting the application on", slog.String("port", port))
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error on server:", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	log.Println("shutting down the server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("error on server shutdown:", err)
		return
	}

	log.Println("shutting down the application")
	if err := app.GracefulShutdown(ctx); err != nil {
		logger.Error("error on application shutdown:", err)
		return
	}

	logger.Info("application stopped")
}
