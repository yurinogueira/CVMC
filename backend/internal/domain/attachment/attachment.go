package attachment

import "time"

type Attachment struct {
	ID            string
	MaintenanceID string
	FileName      string
	OriginalName  string
	MimeType      string
	Size          int64
	Hash          string
	UploadedAt    time.Time
}
