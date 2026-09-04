package maintenance

import "time"

type Attachment struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Size      int64     `json:"size" bson:"size"`
	MimeType  string    `json:"mimeType" bson:"mimeType"`
	DataUrl   string    `json:"dataUrl" bson:"dataUrl"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type Maintenance struct {
	ID          string       `json:"id" bson:"_id"`
	CarID       string       `json:"carId" bson:"carId"`
	Title       string       `json:"title" bson:"title"`
	Description string       `json:"description" bson:"description"`
	Date        time.Time    `json:"date" bson:"date"`
	Mileage     int          `json:"mileage" bson:"mileage"`
	Types       []string     `json:"types,omitempty" bson:"types,omitempty"`
	Cost        *float64     `json:"cost,omitempty" bson:"cost,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty" bson:"attachments,omitempty"`
	CreatedAt   time.Time    `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt" bson:"updatedAt"`
}
