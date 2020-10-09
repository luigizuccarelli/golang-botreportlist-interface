package schema

// Response schema
type Response struct {
	Code    int          `json:"code"`
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Reports []ReportList `json:"reports,omitempty"`
}

// ResponseCount schema
type ResponseCount struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Count   int64  `json:"count"`
}

// StatsResponse schema
type StatsResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Stats   []Stat `json:"stats,omitempty"`
}

// ReportResponse schema
type ReportResponse struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Report  ReportContent `json:"report,omitempty"`
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
	JwtToken string     `json:"jwtToken"`
	Data     ReportList `json:"data,omitempty"`
}

// ReportContent schema
type ReportContent struct {
	Channel             string         `json:"Channel"`
	Affiliate           string         `json:"Affiliate"`
	MessageID           string         `json:"MessageId"`
	EmailBody           string         `json:"EmailBody"`
	EmailSubject        string         `json:"EmailSubject"`
	EmailAdress         string         `json:"EmailAddress"`
	EmailRecipient      string         `json:"EmailRecipient"`
	EmailS3Key          string         `json:"EmailS3Key"`
	Timestamp           int64          `json:"Timestamp"`
	Endpoint            interface{}    `json:"Endpoint"`
	BotProcessingMode   string         `json:"BotProcessingMode"`
	ProcessOutcome      string         `json:"ProcessOutcome"`
	Entities            []string       `json:"Entities"`
	EmailClassification string         `json:"EmailClassification"`
	UserClassification  string         `json:"UserClassification"`
	Success             bool           `json:"Success"`
	CustomerInfo        CustomerDetail `json:"CustomerInfo"`
}

type CustomerDetail struct {
	CustomerNumber  string `json:"customerNumber"`
	ExpirationDate  string `json:"expirationDate"`
	IssuesRemaining int64  `json:"issuesRemaining"`
	CircStatus      string `json:"circStatus"`
	RenewalFlag     string `json:"renewalFlag"`
	ProductFamily   string `json:"productFamily"`
	PubCode         string `json:"pubcode"`
	SubRef          string `json:"subref"`
	Message         string `json:"message"`
}

// Stats schema
type Stat struct {
	Count   float64 `json:"count"`
	Success bool    `json:"success"`
}

// List schema
type ListObject struct {
	ProcessOutcome      string `json:"ProcessOutcome"`
	EmailClassification string `json:"EmailClassification"`
	UserClassification  string `json:"UserClassification"`
	Success             bool   `json:"Success"`
	Timestamp           int64  `json:"Timestamp"`
	AffiliateId         string `json:"AffiliateId"`
}

type ReportList struct {
	Id             string     `json:"id"`
	ServisbotStats ListObject `json:"servisbotstats"`
}
