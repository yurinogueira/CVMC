package car

import "time"

type Car struct {
	ID              string    `json:"id"`
	OwnerID         string    `json:"ownerId"`
	Name            string    `json:"name"`
	Manufacturer    string    `json:"manufacturer"`
	Model           string    `json:"model"`
	YearManufacture int       `json:"yearManufacture"`
	YearModel       int       `json:"yearModel"`
	LastMileage     int       `json:"lastMileage"`
	SharedWith      []string  `json:"sharedWith,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
