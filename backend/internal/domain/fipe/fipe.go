package fipe

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidVehicleType = errors.New("invalid vehicle type: must be cars, motorcycles, or trucks")
	ErrBrandNotFound      = errors.New("brand not found")
	ErrModelNotFound      = errors.New("model not found")
	ErrYearNotFound       = errors.New("year not found")
)

type VehicleType string

const (
	VehicleTypeCars        VehicleType = "cars"
	VehicleTypeMotorcycles VehicleType = "motorcycles"
	VehicleTypeTrucks      VehicleType = "trucks"
)

func ParseVehicleType(s string) (VehicleType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "cars", "carro", "carros", "1":
		return VehicleTypeCars, nil
	case "motorcycles", "motos", "moto", "2":
		return VehicleTypeMotorcycles, nil
	case "trucks", "caminhoes", "caminhão", "caminhao", "3":
		return VehicleTypeTrucks, nil
	default:
		return "", ErrInvalidVehicleType
	}
}

type Brand struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	VehicleType string `json:"vehicleType,omitempty"`
}

type Model struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Year struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type VehicleDetail struct {
	Brand          string  `json:"brand"`
	CodeFipe       string  `json:"codeFipe"`
	Fuel           string  `json:"fuel"`
	FuelAcronym    string  `json:"fuelAcronym,omitempty"`
	Model          string  `json:"model"`
	ModelYear      int     `json:"modelYear"`
	Price          string  `json:"price"`
	PriceValue     float64 `json:"priceValue,omitempty"`
	ReferenceMonth string  `json:"referenceMonth"`
	VehicleType    int     `json:"vehicleType"`
}

type BrandDocument struct {
	ID               string          `bson:"_id,omitempty" json:"id,omitempty"`
	Code             string          `bson:"code" json:"code"`
	Name             string          `bson:"name" json:"name"`
	VehicleType      string          `bson:"vehicleType" json:"vehicleType"`
	ModelsLastSyncAt time.Time       `bson:"modelsLastSyncAt,omitempty" json:"modelsLastSyncAt,omitempty"`
	Models           []ModelDocument `bson:"models,omitempty" json:"models,omitempty"`
	CreatedAt        time.Time       `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt        time.Time       `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type ModelDocument struct {
	Code            string         `bson:"code" json:"code"`
	Name            string         `bson:"name" json:"name"`
	YearsLastSyncAt time.Time      `bson:"yearsLastSyncAt,omitempty" json:"yearsLastSyncAt,omitempty"`
	Years           []YearDocument `bson:"years,omitempty" json:"years,omitempty"`
}

type YearDocument struct {
	Code              string    `bson:"code" json:"code"`
	Name              string    `bson:"name" json:"name"`
	Price             string    `bson:"price,omitempty" json:"price,omitempty"`
	PriceValue        float64   `bson:"priceValue,omitempty" json:"priceValue,omitempty"`
	FIPECode          string    `bson:"fipeCode,omitempty" json:"fipeCode,omitempty"`
	Fuel              string    `bson:"fuel,omitempty" json:"fuel,omitempty"`
	ReferenceMonth    string    `bson:"referenceMonth,omitempty" json:"referenceMonth,omitempty"`
	DetailsLastSyncAt time.Time `bson:"detailsLastSyncAt,omitempty" json:"detailsLastSyncAt,omitempty"`
}
