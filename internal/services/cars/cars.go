package cars

import (
	"boilerplate/internal/models"
	"boilerplate/internal/repository/memory"
	"context"
	"errors"
	"log/slog"

	"go.opentelemetry.io/otel"
)

var ErrCarNotFound = errors.New("car not found")

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

// Find returns a car by id.
func (s *Service) Find(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := tracer.Start(ctx, "Find")
	defer span.End()
	car, err := s.Repository.Find(ctx, id)
	if err != nil && errors.Is(err, memory.ErrCarNotFound) {
		return nil, ErrCarNotFound
	}
	return car, nil
}

// FindAll returns all cars.
func (s *Service) FindAll(ctx context.Context) ([]models.Car, error) {
	ctx, span := tracer.Start(ctx, "FindAll")
	defer span.End()
	return s.Repository.FindAll(ctx)
}

// Create creates a car.
func (s *Service) Create(ctx context.Context, car models.Car) (*models.Car, error) {
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()
	return s.Repository.Create(ctx, car)
}

// Update updates a car.
func (s *Service) Update(ctx context.Context, car models.Car) (*models.Car, error) {
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()
	reponse, err := s.Repository.Update(ctx, car)
	if err != nil && errors.Is(err, memory.ErrCarNotFound) {
		return nil, ErrCarNotFound
	}
	return reponse, nil
}

func New(repo Repository, log *slog.Logger) *Service {
	return &Service{
		Repository: repo,
		log:        log,
	}
}
