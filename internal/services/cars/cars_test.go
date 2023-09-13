package cars_test

import (
	"boilerplate/internal/models"
	"boilerplate/internal/repository/memory"
	"boilerplate/internal/services/cars"
	"boilerplate/internal/services/cars/mock"
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestService_Find(t *testing.T) {
	now := time.Now()
	car := models.Car{
		ID:        "1",
		Model:     "hb20",
		Category:  "Hatch",
		Year:      2023,
		Price:     100000,
		Color:     "White",
		Make:      "Ford",
		Mileage:   1000,
		Package:   "SE",
		UpdatedAt: now,
		CreatedAt: now.Add(-time.Hour),
	}

	t.Parallel()
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		mock    func(mock *mock.MockRepository)
		want    *models.Car
		wantErr bool
		err     error
	}{
		{
			name: "Should return an error when car is not found",
			args: args{id: "1"},
			mock: func(mock *mock.MockRepository) {
				mock.EXPECT().Find(gomock.Any(), gomock.Eq("1")).Return(nil, memory.ErrCarNotFound)
			},
			err:     cars.ErrCarNotFound,
			wantErr: true,
		},
		{
			name: "Should return a car when car is found",
			args: args{id: "1"},
			mock: func(mock *mock.MockRepository) {
				mock.EXPECT().Find(gomock.Any(), gomock.Eq("1")).Return(&car, nil)
			},
			want: &car,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mock := mock.NewMockRepository(ctrl)
			tt.mock(mock)

			s := cars.New(mock, slog.Default())

			got, err := s.Find(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.err) {
				t.Errorf("Service.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_FindAll(t *testing.T) {
	t.Parallel()
	now := time.Now()
	car := models.Car{
		ID:        "1",
		Model:     "hb20",
		Category:  "Hatch",
		Year:      2023,
		Price:     100000,
		Color:     "White",
		Make:      "Ford",
		Mileage:   1000,
		Package:   "SE",
		UpdatedAt: now,
		CreatedAt: now.Add(-time.Hour),
	}

	tests := []struct {
		name    string
		mock    func(mock *mock.MockRepository)
		want    []models.Car
		wantErr bool
	}{
		{
			name: "Should return an epmty list of cars",
			mock: func(mock *mock.MockRepository) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]models.Car{}, nil)
			},
			want: []models.Car{},
		},
		{
			name: "Should return a list of cars",
			mock: func(mock *mock.MockRepository) {
				mock.EXPECT().FindAll(gomock.Any()).Return([]models.Car{
					car,
				}, nil)
			},
			want: []models.Car{
				car,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mock := mock.NewMockRepository(ctrl)
			tt.mock(mock)

			s := cars.New(mock, slog.Default())

			got, err := s.FindAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	t.Parallel()
	now := time.Now()
	car := models.Car{
		ID:        "1",
		Model:     "hb20",
		Category:  "Hatch",
		Year:      2023,
		Price:     100000,
		Color:     "White",
		Make:      "Ford",
		Mileage:   1000,
		Package:   "SE",
		UpdatedAt: now,
		CreatedAt: now,
	}

	type args struct {
		car models.Car
	}
	tests := []struct {
		name    string
		mock    func(mock *mock.MockRepository)
		args    args
		want    *models.Car
		wantErr bool
	}{
		{
			name: "Should create an car",
			args: args{car: car},
			mock: func(mock *mock.MockRepository) {
				mock.EXPECT().Create(gomock.Any(), gomock.Eq(car)).Return(&car, nil)
			},
			want: &car,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mock := mock.NewMockRepository(ctrl)
			tt.mock(mock)

			s := cars.New(mock, slog.Default())

			got, err := s.Create(context.Background(), tt.args.car)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()
	now := time.Now()
	car := models.Car{
		ID:        "1",
		Model:     "hb20",
		Category:  "Hatch",
		Year:      2023,
		Price:     100000,
		Color:     "White",
		Make:      "Ford",
		Mileage:   1000,
		Package:   "SE",
		UpdatedAt: now,
		CreatedAt: now,
	}

	type args struct {
		car models.Car
	}
	tests := []struct {
		name    string
		args    args
		mock    func(mock *mock.MockRepository)
		want    *models.Car
		err     error
		wantErr bool
	}{
		{
			name: "Should update an car",
			args: args{car: car},
			mock: func(mock *mock.MockRepository) {
				mock.EXPECT().Update(gomock.Any(), gomock.Eq(car)).Return(&car, nil)
			},
			want: &car,
		},
		{
			name: "Should return an error when car is not found",
			args: args{car: car},
			mock: func(mock *mock.MockRepository) {
				mock.EXPECT().Update(gomock.Any(), gomock.Eq(car)).Return(nil, memory.ErrCarNotFound)
			},
			err:     cars.ErrCarNotFound,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mock := mock.NewMockRepository(ctrl)
			tt.mock(mock)

			s := cars.New(mock, slog.Default())

			got, err := s.Update(context.Background(), tt.args.car)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !errors.Is(err, tt.err) {
				t.Errorf("Service.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
