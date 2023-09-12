package models

import "time"

type Car struct {
	ID       string  `json:"id"`
	Category string  `json:"category"`
	Color    string  `json:"color"`
	Make     string  `json:"make"`
	Mileage  int     `json:"mileage"`
	Model    string  `json:"model"`
	Package  string  `json:"package"`
	Price    float64 `json:"price"`
	Year     int     `json:"year"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
