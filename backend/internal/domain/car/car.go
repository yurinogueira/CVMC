package car

import "time"

type Car struct {
	ID                 string
	OwnerID            string
	Name               string
	Manufacturer       string
	Model              string
	YearManufacture    int
	YearModel          int
	LastMileage        int
	SharedWith         []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
