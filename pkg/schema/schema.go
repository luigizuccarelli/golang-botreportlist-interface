package schema

import "time"

// Response schema
type Response struct {
	Code         int              `json:"name"`
	StatusCode   string           `json:"statuscode"`
	Status       string           `json:"status"`
	Message      string           `json:"message"`
	Payload      []S3FileMetaData `json:"payload"`
	EmailContent string           `json:"email,omitempty"`
}

type S3FileMetaData struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"lastmodified"`
	Size         int64     `json:"size"`
	StorageClass string    `json:"class"`
}
