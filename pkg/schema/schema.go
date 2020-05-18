package schema

// Response schema
type Response struct {
	Name       string           `json:"name"`
	StatusCode string           `json:"statuscode"`
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	Payload    *SchemaInterface `json:"payload"`
}

// Token Schema
type TokenDetail struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type EmailProfile struct {
	CviNbr         string `json:"cviNbr"`
	CustomerNumber string `json:"customerNumber"`
	Role           string `json:"role"`
	Temp           bool   `json:"temp"`
	Password       string `json:"password"`
	ID             struct {
		UserName   string `json:"userName"`
		PortalCode struct {
			AuthGroup string `json:"authGroup"`
		} `json:"portalCode"`
		AuthType string `json:"authType"`
	} `json:"id"`
	DenyAccess string `json:"denyAccess"`
	AuthStatus string `json:"authStatus"`
}

type Email struct {
	EmailAddress string `json: "emailAddress"`
	Temp         bool   `json: "temp"`
}

type PostalAddress struct {
	BirthDate    string `json: "birthDate"`
	City         string `json: "city"`
	CountryCode  string `json: "countryCode"`
	EmailAddress Email
	FirstName    string `json: "firstName"`
	LastName     string `json: "lastName"`
	PhoneNumber  string `json: "phoneNumber"`
	PostalCode   string `json: "postalCode"`
	State        string `json: "state"`
	Street       string `json: "street"`
	Street2      string `json: "street2"`
	Street3      string `json: "street3"`
	Title        string `json: "title"`
}

type SubscriptionsOrders struct {
	AccessMaintenanceOrders []interface{} `json: "accessMaintenanceOrders,omitempty"`
	ProductOrders           []interface{} `json: "productOrders,omitempty"`
	Subscriptions           []Subscription
}

// Subscription - the main data
// This schema will change for each microservice
type Subscription struct {
	CircStatus          string `json:"circStatus"`
	DeliveryCode        string `json:"deliveryCode"`
	PromoCode           string `json:"promoCode"`
	StartDate           string `json:"startDate"`
	FinalExpirationDate string `json:"finalExpirationDate"`
	IssuesRemaining     int64  `json:"issuesRemaining"`
	LastIssue           string `json:"lastIssue"`
	Id                  ID     `json:"id"`
}

type ID struct {
	CustomerNumber string   `json:"customerNumber"`
	SubRef         string   `json:"subRef"`
	Item           ItemData `json:"item"`
}

type ItemData struct {
	Format              string `json:"format"`
	ItemDemoGraphic     string `json: "itemDemographic1"`
	ItemDescription     string `json: "itemDescription"`
	ItemNumber          string `json: "itemNumber"`
	ItemType            string `json: "itemType"`
	OwningOrg           string `json: "owningOrg"`
	PackageFlag         string `json: "packageFlag"`
	ProductFamily       string `json: "productFamily"`
	SaleableFlag        string `json: "saleableFlag"`
	SalesTaxProductCode string `json: "salesTaxProductCode"`
	ServiceCode         string `json: "serviceCode"`
	TaxCommodityCode    string `json: "taxCommodityCode"`
}

// ShcemaInterface - acts as an interface wrapper for our profile schema
// All the go microservices will using this schema
type SchemaInterface struct {
	ID                     int64               `json:"_id,omitempty"`
	LastUpdate             int64               `json:"lastupdate,omitempty"`
	MetaInfo               string              `json:"metainfo,omitempty"`
	Accounts               []EmailProfile      `json:"accounts,omitempty"`
	SubscriptionsAndOrders SubscriptionsOrders `json:"subscriptions,omitempty"`
	PostalAddresses        []PostalAddress     `json:"postaladdresses,omitempty"`
	EmailAddresses         []Email             `json:"emailaddresses,omitempty"`
}
