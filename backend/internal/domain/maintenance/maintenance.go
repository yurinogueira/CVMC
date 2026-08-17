package maintenance

import "time"

type Maintenance struct {
	ID          string
	CarID       string
	Title       string
	Description string
	Date        time.Time
	Mileage     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
