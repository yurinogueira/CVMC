package fuel

import "time"

type Fueling struct {
	ID            string    `json:"id" bson:"_id"`
	CarID         string    `json:"carId" bson:"carId"`
	Date          time.Time `json:"date" bson:"date"`
	FuelType      string    `json:"fuelType" bson:"fuelType"`
	Liters        float64   `json:"liters" bson:"liters"`
	PricePerLiter float64   `json:"pricePerLiter" bson:"pricePerLiter"`
	TotalCost     float64   `json:"totalCost" bson:"totalCost"`
	IsFullTank    bool      `json:"isFullTank" bson:"isFullTank"`
	GasStation    string    `json:"gasStation,omitempty" bson:"gasStation,omitempty"`
	Notes         string    `json:"notes,omitempty" bson:"notes,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}
