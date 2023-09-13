package controler

import (
	"boilerplate/internal/models"
	"boilerplate/internal/services/cars"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
)

type CarsControler struct {
	Service *cars.Service
	log     *slog.Logger
}

// HandleGetCars returns all cars.
func (c *CarsControler) HandleGetCars(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(pkgName).Start(r.Context(), "HandleGetCars")
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	cars, err := c.Service.FindAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(cars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(js); err != nil {
		c.log.Error("failed to write response: ", err)
		return
	}
}

// HandleGetCarByID returns a car by id.
func (c *CarsControler) HandleGetCarByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(pkgName).Start(r.Context(), "HandleGetCarByID")
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	pk := ctx.Value(httprouter.ParamsKey)

	ps, ok := pk.(httprouter.Params)
	if !ok {
		http.Error(w, "invalid code", http.StatusBadRequest)
		return
	}

	id := ps.ByName("id")
	if id == "" {
		http.Error(w, "invalid code", http.StatusBadRequest)
	}

	car, err := c.Service.Find(ctx, id)
	if err != nil {
		if errors.Is(err, cars.ErrCarNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(js); err != nil {
		c.log.Error("failed to write response: ", err)
		return
	}
}

// HandleCreateCar creates a car.
func (c *CarsControler) HandleCreateCar(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(pkgName).Start(r.Context(), "HandleCreateCar")
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := c.Service.Create(ctx, car)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(js); err != nil {
		c.log.Error("failed to write response: ", err)
		return
	}
}

// HandleUpdateCar updates a car.
func (c *CarsControler) HandleUpdateCar(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(pkgName).Start(r.Context(), "HandleUpdateCar")
	defer span.End()
	w.Header().Set("Content-Type", "application/json")
	pk := ctx.Value(httprouter.ParamsKey)

	ps, ok := pk.(httprouter.Params)
	if !ok {
		http.Error(w, "invalid code", http.StatusBadRequest)
		return
	}

	id := ps.ByName("id")
	if id == "" {
		http.Error(w, "invalid code", http.StatusBadRequest)
	}

	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	car.ID = id
	response, err := c.Service.Update(ctx, car)
	if err != nil {
		if errors.Is(err, cars.ErrCarNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(js); err != nil {
		c.log.Error("failed to write response: ", err)
		return
	}
}

func NewCarsControler(service *cars.Service, log *slog.Logger) *CarsControler {
	return &CarsControler{
		Service: service,
		log:     log,
	}
}
