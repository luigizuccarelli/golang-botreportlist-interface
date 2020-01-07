package main

// Analytics schema that forwards the json data payload to our backend analytics system (Couchbase)
type Analytics struct {
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
