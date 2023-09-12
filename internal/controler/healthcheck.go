package controler

import (
	"boilerplate/internal/services/healthcheck"
	"fmt"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const pkgName = "controler"

type HealthCheckControler struct {
	alive *healthcheck.Alive
	log   *slog.Logger
}

func (ctrl *HealthCheckControler) HandleHeathCheck(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer(pkgName).Start(r.Context(), "HandleHeathCheck")
	defer span.End()
	if _, err := w.Write([]byte(healthcheck.OK)); err != nil {
		ctrl.log.Error("failed to write response: ", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ctrl *HealthCheckControler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(pkgName).Start(r.Context(), "HandleReadiness")
	defer span.End()

	if err := ctrl.alive.Readiness(ctx); err != nil {
		ctrl.log.Error("Readiness failed: ", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if _, err = w.Write([]byte(fmt.Sprintf("%s", err))); err != nil {
			ctrl.log.Error("failed to write response: ", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(healthcheck.OK)); err != nil {
		ctrl.log.Error("failed to write response: ", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewHealthCheckControler(log *slog.Logger, alive *healthcheck.Alive) *HealthCheckControler {
	return &HealthCheckControler{
		alive: alive,
		log:   log,
	}
}
