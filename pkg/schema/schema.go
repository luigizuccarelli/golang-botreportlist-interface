package schema

import "time"

// Response schema
type Response struct {
	Code       int              `json:"name"`
	StatusCode string           `json:"statuscode"`
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	Payload    []S3FileMetaData `json:"payload,omitempty"`
	Email      string           `json:"email,omitempty"`
	Report     *ReportContent   `json:"report,omitempty"`
}

// GenericSchema - used in the GenericHandler (complex data object)
type GenericSchema struct {
	Token   string
	Creds   *Credentials
	Request *ServisBOTRequest
}

// Token Schema
type TokenDetail struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// Credentials (from JWT)
type Credentials struct {
	User           string `json:"user"`
	Password       string `json:"password"`
	CustomerNumber string `json:"customerNumber"`
}

type ServisBOTRequest struct {
	Email          string `json:"email"`
	JwtToken       string `json:"jwtToken"`
	Subscription   string `json:"subscription"`
	Reason         string `json:"reason"`
	CustomerNumber string `json:"customerNumber"`
	SubRef         string `json:"subref,omitempty"`
	PubCode        string `json:"pubcode,omitempty"`
	UniqueId       string `json:"uniqueid,omitempty"`
	PhoneNumber    string `json:"phonenumber,omitempty"`
	RenewalFlag    string `json:"renewalFlag,omitempty"`
	Subject        string `json:"subject,omitempty"`
	Note           string `json:"note,omitempty"`
	Data           string `json:"data,omitempty"`
}

type S3FileMetaData struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"lastmodified"`
	Size         int64     `json:"size"`
	StorageClass string    `json:"class"`
}

type ReportContent struct {
	Channel             string      `json:"Channel"`
	Affiliate           string      `json:"Affiliate"`
	MessageID           string      `json:"MessageId"`
	EmailBody           string      `json:"EmailBody"`
	EmailSubject        string      `json:"EmailSubject"`
	EmailS3Key          string      `json:"EmailS3Key"`
	Timestamp           int64       `json:"Timestamp"`
	Endpoint            interface{} `json:"Endpoint"`
	BotProcessingMode   string      `json:"BotProcessingMode"`
	ProcessOutcome      string      `json:"ProcessOutcome"`
	Entities            []string    `json:"Entities"`
	EmailClassification []string    `json:"EmailClassification"`
}
