package healthcheck

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
)

const pkgName = "healthcheck"
const OK = "OK"

var ErrDepencenciesFailed = errors.New("depencencies failed")

//go:generate mockgen -source=healthcheck.go -destination=mock/healthcheck.go -package=mock
type Dependency interface {
	Healthcheck(context.Context) error
}

type Alive struct {
	depencencies []Dependency
}

func (alive *Alive) Readiness(ctx context.Context) error {
	ctx, span := otel.Tracer(pkgName).Start(ctx, "Readiness")
	defer span.End()

	g, ctx := errgroup.WithContext(ctx)

	for _, depencency := range alive.depencencies {
		depencency := depencency
		g.Go(func() error {
			return depencency.Healthcheck(ctx)
		})
	}

	if err := g.Wait(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.Join(ErrDepencenciesFailed, err)
	}

	return nil
}

func New(depencencies ...Dependency) *Alive {
	return &Alive{
		depencencies: depencencies,
	}
}
