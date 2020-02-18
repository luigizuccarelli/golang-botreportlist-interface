package schema

import (
	"time"
)

// Response schema
type Response struct {
	Name       string      `json:"name"`
	StatusCode string      `json:"statuscode"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload"`
}

// Analytics schema that forwards the json data payload to our backend analytics system (Couchbase)
type Analytics struct {
	Id         string         `json:"id"`
	TrackingId string         `json:"trackingid"`
	To         PageDetail     `json:"to"`
	From       PageDetail     `json:"from"`
	Location   LocationDetail `json:"location"`
	Currency   CurrencyDetail `json:"currency"`
	Event      EventDetail    `json:"event"`
	Campaign   string         `json:"utm_campaign"`
	Affiliate  string         `json:"utm_affiliate"`
	Medium     string         `json:"utm_medium"`
	Source     string         `json:"utm_source"`
	Conent     string         `json:"utm_content"`
	Timestamp  int64          `json:"timestamp"`
	Platform   PlatformDetail `json:"platform"`
	Creative   CreativeDetail `json:"creative"`
	Effort     EffortDetail   `json:"effort"`
	Journey    JourneyDetail  `json:"journey"`
	Product    string         `json:"product"`
}

type PageDetail struct {
	Url      string `json:"url"`
	PageName string `json:"pagename"`
	PageType string `json:"pagetype"`
}

// LocationDetail schema - used in the Analytics schema
type LocationDetail struct {
	Ip      string        `json:"ip"`
	Carrier string        `json:"carrier"`
	Country CountryDetail `json:"country"`
}

// CountryDetail schema
type CountryDetail struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Capitial string `json:"capital"`
}

// CurrencyDetail schema - used in the Analytics schema
type CurrencyDetail struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// EventDetail schema
type EventDetail struct {
	Type       string `json:"type"`
	TimeonPage int    `json:"timeonpage"`
}

// PlatformDetail
type PlatformDetail struct {
	AppcodeName string `json:"appCodeName"`
	AppName     string `json:"appName"`
	AppVersion  string `json:"appVersion"`
	Language    string `json:"language"`
	Os          string `json:"os"`
	Product     string `json:"product"`
	ProductSub  string `json:"productSub"`
	UserAgent   string `json:"userAgent"`
	Vendor      string `json:"vendor"`
}

type CreativeDetail struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type EffortDetail struct {
	AcquisitionMethod        string `json:"acquisition_method"`
	AdvantageCampaignCode    string `json:"advantage_campaign_code"`
	AdvantageDescription     string `json:"advantage_description"`
	AdvertisementName        string `json:"advertisement_name"`
	Campaign                 string `json:"campaign"`
	Domain                   string `json:"domain"`
	Date                     string `json:"date"`
	EffortDestination        string `json:"effor_destination"`
	Id                       string `json:"id"`
	Type                     string `json:"type"`
	Journey                  string `json:"journey"`
	Promocode                string `json:"promocode"`
	WhatAreYouPomoting       string `json:"what_are_you_promoting"`
	WhereIsTheMarketingGoing string `json:"where_is_the_marketing_going"`
}

type JourneyDetail struct {
	CreativeSequence string `json:"creative_sequence"`
	Name             string `json:"name"`
	Status           string `json:"status"`
}

type SegmentIO struct {
	Id          string `json:"id,omitempy"`
	AnonymousID string `json:"anonymousId"`
	Context     struct {
		Campaign struct {
			Content string `json:"content"`
			Name    string `json:"name"`
			Source  string `json:"source"`
		} `json:"campaign"`
		IP      string `json:"ip"`
		Library struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"library"`
		Locale string `json:"locale"`
		Page   struct {
			Path     string `json:"path"`
			Referrer string `json:"referrer"`
			Search   string `json:"search"`
			Title    string `json:"title"`
			URL      string `json:"url"`
		} `json:"page"`
		UserAgent string `json:"userAgent"`
	} `json:"context"`
	Event        string `json:"event"`
	Integrations struct {
	} `json:"integrations"`
	MessageID         string    `json:"messageId"`
	OriginalTimestamp time.Time `json:"originalTimestamp"`
	Properties        struct {
		IrisPlusData struct {
			CREATIVEName                   string `json:"CREATIVE.name"`
			CREATIVEStatus                 string `json:"CREATIVE.status"`
			EFFORTAdvertisementName        string `json:"EFFORT.advertisement_name"`
			EFFORTCampaign                 string `json:"EFFORT.campaign"`
			EFFORTDate                     string `json:"EFFORT.date"`
			EFFORTDomain                   string `json:"EFFORT.domain"`
			EFFORTEffortDestination        string `json:"EFFORT.effort_destination"`
			EFFORTID                       string `json:"EFFORT.id"`
			EFFORTJourneyName              string `json:"EFFORT.journey_name"`
			EFFORTPromocode                string `json:"EFFORT.promocode"`
			EFFORTType                     string `json:"EFFORT.type"`
			EFFORTWhatAreYouPromoting      string `json:"EFFORT.what_are_you_promoting"`
			EFFORTWhereIsTheMarketingGoing string `json:"EFFORT.where_is_the_marketing_going"`
			FORMINFOAddress                string `json:"FORM_INFO.address"`
			FORMINFOAddress2               string `json:"FORM_INFO.address2"`
			FORMINFOAddress3               string `json:"FORM_INFO.address3"`
			FORMINFOCity                   string `json:"FORM_INFO.city"`
			FORMINFOCompanyName            string `json:"FORM_INFO.companyName"`
			FORMINFOCountryCode            string `json:"FORM_INFO.countryCode"`
			FORMINFOCountryName            string `json:"FORM_INFO.countryName"`
			FORMINFOEmail                  string `json:"FORM_INFO.email"`
			FORMINFOFaxNumber              string `json:"FORM_INFO.faxNumber"`
			FORMINFOFirstName              string `json:"FORM_INFO.firstName"`
			FORMINFOLastName               string `json:"FORM_INFO.lastName"`
			FORMINFOPhoneNumber            string `json:"FORM_INFO.phoneNumber"`
			FORMINFOPhoneNumber2           string `json:"FORM_INFO.phoneNumber2"`
			FORMINFOPhoneNumber3           string `json:"FORM_INFO.phoneNumber3"`
			FORMINFOPostalCode             string `json:"FORM_INFO.postalCode"`
			FORMINFOStateCode              string `json:"FORM_INFO.stateCode"`
			FORMINFOStateName              string `json:"FORM_INFO.stateName"`
			FORMINFOSuffix                 string `json:"FORM_INFO.suffix"`
			FORMINFOTitle                  string `json:"FORM_INFO.title"`
			JOURNEYCreativeSequence        string `json:"JOURNEY.creative_sequence"`
			JOURNEYName                    string `json:"JOURNEY.name"`
			JOURNEYStatus                  string `json:"JOURNEY.status"`
		} `json:"irisPlusData"`
		Type        string `json:"type"`
		UtmVariable struct {
			Pagetype    string `json:"pagetype"`
			PromoCode   string `json:"promoCode"`
			UtmCampaign string `json:"utm_campaign"`
			UtmContent  string `json:"utm_content"`
			UtmSource   string `json:"utm_source"`
		} `json:"utm_variable"`
		Value int `json:"value"`
	} `json:"properties"`
	ReceivedAt time.Time `json:"receivedAt"`
	SentAt     time.Time `json:"sentAt"`
	Timestamp  time.Time `json:"timestamp"`
	Type       string    `json:"type"`
	UserID     string    `json:"userId"`
	Version    int       `json:"version"`
}
