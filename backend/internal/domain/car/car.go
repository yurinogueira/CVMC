package car

import "time"

type MaintenanceIntervalConfig struct {
	IntervalKm     int `json:"intervalKm" bson:"intervalKm"`
	IntervalMonths int `json:"intervalMonths" bson:"intervalMonths"`
}

type Car struct {
	ID                   string                               `json:"id"`
	OwnerID              string                               `json:"ownerId"`
	Name                 string                               `json:"name"`
	Manufacturer         string                               `json:"manufacturer"`
	Model                string                               `json:"model"`
	YearManufacture      int                                  `json:"yearManufacture"`
	YearModel            int                                  `json:"yearModel"`
	InitialMileage       int                                  `json:"initialMileage,omitempty"`
	LastMileage          int                                  `json:"lastMileage"`
	VehicleType          string                               `json:"vehicleType,omitempty"`
	ImageUrl             string                               `json:"imageUrl,omitempty"`
	FIPECode             string                               `json:"fipeCode,omitempty"`
	FIPEPrice            string                               `json:"fipePrice,omitempty"`
	Fuel                 string                               `json:"fuel,omitempty"`
	SharedWith           []string                             `json:"sharedWith,omitempty"`
	MaintenanceIntervals map[string]MaintenanceIntervalConfig `json:"maintenanceIntervals,omitempty"`
	MileageUpdatedAt     time.Time                            `json:"mileageUpdatedAt,omitempty"`
	CreatedAt            time.Time                            `json:"createdAt"`
	UpdatedAt            time.Time                            `json:"updatedAt"`
}

// EffectiveMileageUpdatedAt returns the timestamp of the last mileage update,
// falling back to UpdatedAt and then CreatedAt for legacy records without it.
func (c Car) EffectiveMileageUpdatedAt() time.Time {
	if !c.MileageUpdatedAt.IsZero() {
		return c.MileageUpdatedAt
	}
	if !c.UpdatedAt.IsZero() {
		return c.UpdatedAt
	}
	return c.CreatedAt
}
