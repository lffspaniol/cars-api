package memory

import (
	"boilerplate/internal/models"
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRepository_Find(t *testing.T) {
	t.Parallel()
	car := models.Car{
		ID:       "1",
		Model:    "Fiesta",
		Category: "Hatch",
		Year:     2020,
		Price:    100000,
		Color:    "White",
		Make:     "Ford",
		Mileage:  1000,
		Package:  "SE",
	}

	type fields struct {
		values map[string]models.Car
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Car
		wantErr bool
		err     error
	}{
		{
			name: "should return an error when car is not found",
			fields: fields{
				values: map[string]models.Car{},
			},
			args:    args{id: "1"},
			wantErr: true,
			err:     ErrCarNotFound,
		},
		{
			name: "should an car when car is found",
			fields: fields{
				values: map[string]models.Car{
					"1": car,
				},
			},
			args: args{id: "1"},
			want: &car,
			err:  nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Repository{
				values: tt.fields.values,
			}
			got, err := r.Find(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && errors.Is(tt.err, err) {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_FindAll(t *testing.T) {
	t.Parallel()
	now := time.Now()
	car1 := models.Car{
		ID:        "1",
		Model:     "Fiesta",
		Category:  "Hatch",
		Year:      2020,
		Price:     100000,
		Color:     "White",
		Make:      "Ford",
		Mileage:   1000,
		Package:   "SE",
		CreatedAt: now,
		UpdatedAt: now,
	}

	car2 := models.Car{
		ID:        "2",
		Model:     "HB20",
		Category:  "Hatch",
		Year:      2022,
		Price:     1000000,
		Color:     "silver",
		Make:      "Hyundai",
		Mileage:   10000,
		Package:   "Comfort Plus",
		CreatedAt: now.Add(time.Hour),
		UpdatedAt: now.Add(time.Hour),
	}

	type fields struct {
		values map[string]models.Car
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Car
		wantErr bool
	}{
		{
			name: "should return an emptry slice when there is no cars",
			fields: fields{
				values: map[string]models.Car{},
			},
			want:    []models.Car{},
			wantErr: false,
		},
		{
			name: "should return a slice with all cars",
			fields: fields{
				values: map[string]models.Car{
					"1": car1,
					"2": car2,
				},
			},
			want: []models.Car{car2, car1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Repository{
				values: tt.fields.values,
			}
			got, err := r.FindAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Create(t *testing.T) {
	now := time.Now()
	id := "1"
	idHelper = func() string {
		return id
	}
	car := models.Car{
		ID:        "1",
		Model:     "Fiesta",
		Category:  "Hatch",
		Year:      2020,
		Price:     100000,
		Color:     "White",
		Make:      "Ford",
		Mileage:   1000,
		Package:   "SE",
		CreatedAt: now,
		UpdatedAt: now,
	}
	t.Parallel()

	timeHelper = func() time.Time {
		return now
	}

	type fields struct {
		values map[string]models.Car
	}
	type args struct {
		car models.Car
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Car
		wantErr bool
	}{
		{
			name:   "should create a car",
			fields: fields{values: map[string]models.Car{}},
			args: args{
				car: car,
			},
			want: &car,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Repository{
				values: tt.fields.values,
			}
			got, err := r.Create(context.Background(), tt.args.car)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	now := time.Now()
	car := models.Car{
		ID:        "1",
		Model:     "Fiesta",
		Category:  "Hatch",
		Year:      2020,
		Price:     100000,
		Color:     "White",
		Make:      "Ford",
		Mileage:   1000,
		Package:   "SE",
		CreatedAt: now.Add(-time.Hour),
		UpdatedAt: now.Add(-time.Hour),
	}

	carUpdate := models.Car{
		ID:        "1",
		Model:     "hb20",
		Category:  "Hatch",
		Year:      2023,
		Price:     100000,
		Color:     "White",
		Make:      "Hyundai",
		Mileage:   1000,
		Package:   "SE",
		UpdatedAt: now,
		CreatedAt: now.Add(-time.Hour),
	}
	timeHelper = func() time.Time {
		return now
	}

	type fields struct {
		values map[string]models.Car
	}
	type args struct {
		car models.Car
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Car
		wantErr bool
	}{
		{
			name: "should update a car",
			fields: fields{
				values: map[string]models.Car{
					"1": car,
				},
			},
			args: args{
				car: carUpdate,
			},
			want: &carUpdate,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Repository{
				values: tt.fields.values,
			}
			got, err := r.Update(context.Background(), tt.args.car)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
