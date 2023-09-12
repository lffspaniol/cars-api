package controler

import (
	"boilerplate/internal/models"
	"boilerplate/internal/services/cars"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CarsControler struct {
	Service *cars.Service
	log     *slog.Logger
}

func (c *CarsControler) HandleGetCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cars, err := c.Service.FindAll(r.Context())
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

func (c *CarsControler) HandleGetCarByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
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

func (c *CarsControler) HandleCreateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := c.Service.Create(r.Context(), car)
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

func (c *CarsControler) HandleUpdateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
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
	response, err := c.Service.Update(r.Context(), car)
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

func NewCarsControler(service *cars.Service, log *slog.Logger) *CarsControler {
	return &CarsControler{
		Service: service,
		log:     log,
	}
}
