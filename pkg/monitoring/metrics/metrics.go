package metrics

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"

	"log/slog"
)

// Metrics starts the metrics exporter.
func Metrics(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource, logger *slog.Logger) (func(), error) {
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
		otlpmetricgrpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}

	read := metric.NewPeriodicReader(metricExporter, metric.WithInterval(1*time.Second))

	provider := metric.NewMeterProvider(metric.WithResource(res), metric.WithReader(read))

	otel.SetMeterProvider(provider)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		if shutErr := provider.Shutdown(ctx); shutErr != nil {
			logger.Error("failed to shutdown metrics provider", err)
		}
	}, nil
}
