package container

import (
	"boilerplate/internal/controler"
	"boilerplate/internal/repository/memory"
	"boilerplate/internal/services/cars"
	"boilerplate/internal/services/healthcheck"
	"boilerplate/pkg/monitoring"
	"boilerplate/pkg/monitoring/tracer"
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/spf13/viper"
)

type shutdownFunc func(context.Context) error

type Application struct {
	*controler.HealthCheckControler
	*controler.CarsControler

	Log       *slog.Logger
	shutdowns []shutdownFunc
}

func (app *Application) GracefulShutdown(ctx context.Context) error {
	var err error
	for _, shutdown := range app.shutdowns {
		if err = shutdown(ctx); err != nil {
			err = errors.Join(err, err)
		}
	}
	return err
}

func loadEnvs() {
	viper.SetConfigFile("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Print("using config from env vars instead config file", err)
		} else {
			log.Print("failed to load config using viper", err)
		}
	}

	viper.AutomaticEnv()
}

func startTelemetry(ctx context.Context, logger *slog.Logger) (shutdownFunc, error) {
	res, err := monitoring.Setup()
	if err != nil {
		return nil, err
	}

	con, err := monitoring.Connect(ctx)
	if err != nil {
		if errors.Is(err, monitoring.ErrSDKDisabled) {
			return func(ctx context.Context) error {
				return nil
			}, nil
		}

		return nil, err
	}
	tracer, err := tracer.NewTracer(ctx, con, res, logger)
	if err != nil {
		return nil, err
	}

	return shutdownFunc(tracer), nil
}

func NewApplication(ctx context.Context, logger *slog.Logger) *Application {
	loadEnvs()

	stopTelemetry, err := startTelemetry(ctx, logger)
	if err != nil {
		logger.Error("failed to start telemetry: ", err)
	}

	repositoy := memory.New()
	alive := healthcheck.New(repositoy)
	carsService := cars.New(repositoy, logger)
	healthCheckControler := controler.NewHealthCheckControler(logger, alive)
	carsControler := controler.NewCarsControler(carsService, logger)

	return &Application{
		HealthCheckControler: healthCheckControler,
		CarsControler:        carsControler,
		Log:                  logger,
		shutdowns: []shutdownFunc{
			stopTelemetry,
		},
	}
}
