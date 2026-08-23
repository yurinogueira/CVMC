package maintenance

import "time"

type Maintenance struct {
	ID          string    `json:"id"`
	CarID       string    `json:"carId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Mileage     int       `json:"mileage"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
