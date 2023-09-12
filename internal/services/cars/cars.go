package cars

import (
	"boilerplate/internal/models"
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("cars")

//go:generate mockgen -destination=mock/mock_repository.go -package=mock boilerplate/internal/services/cars Repository
type Repository interface {
	Find(ctx context.Context, id string) (*models.Car, error)
	FindAll(ctx context.Context) ([]models.Car, error)
	Create(ctx context.Context, car models.Car) (*models.Car, error)
	Update(ctx context.Context, car models.Car) (*models.Car, error)
}

type Service struct {
	Repository Repository
	log        *slog.Logger
}

func (s *Service) Find(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := tracer.Start(ctx, "Find")
	defer span.End()
	return s.Repository.Find(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]models.Car, error) {
	ctx, span := tracer.Start(ctx, "FindAll")
	defer span.End()
	return s.Repository.FindAll(ctx)
}

func (s *Service) Create(ctx context.Context, car models.Car) (*models.Car, error) {
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()
	return s.Repository.Create(ctx, car)
}

func (s *Service) Update(ctx context.Context, car models.Car) (*models.Car, error) {
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()
	return s.Repository.Update(ctx, car)
}

func New(repo Repository, log *slog.Logger) *Service {
	return &Service{
		Repository: repo,
		log:        log,
	}
}
