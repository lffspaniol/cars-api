package memory

import (
	"boilerplate/internal/models"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var ErrCarNotFound = fmt.Errorf("car not found")
var timeHelper = time.Now
var idHelper = uuid.NewString
var tracer = otel.Tracer("memory")

type Repository struct {
	values map[string]models.Car
}

func New() *Repository {
	return &Repository{
		values: make(map[string]models.Car),
	}
}

func (r *Repository) Find(ctx context.Context, id string) (*models.Car, error) {
	_, span := tracer.Start(ctx, "Find")
	defer span.End()
	car, ok := r.values[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrCarNotFound, id)
	}
	return &car, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]models.Car, error) {
	_, span := tracer.Start(ctx, "FindAll")
	defer span.End()
	cars := []models.Car{}
	for _, v := range r.values {
		cars = append(cars, v)
	}
	sort.Slice(cars, func(i, j int) bool {
		return cars[i].CreatedAt.After(cars[j].CreatedAt)
	})

	return cars, nil
}

func (r *Repository) Create(ctx context.Context, car models.Car) (*models.Car, error) {
	_, span := tracer.Start(ctx, "Create")
	defer span.End()
	id := idHelper()
	now := timeHelper()
	car.CreatedAt = now
	car.UpdatedAt = now
	r.values[id] = car
	car.ID = id
	return &car, nil
}

func (r *Repository) Update(ctx context.Context, car models.Car) (*models.Car, error) {
	_, span := tracer.Start(ctx, "Update")
	defer span.End()
	car.UpdatedAt = timeHelper()
	r.values[car.ID] = car
	return &car, nil
}

func (Repository) Healthcheck(context.Context) error {
	return nil
}
