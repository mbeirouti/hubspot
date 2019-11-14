package models

import (
	"time"
)

const (
	ISO8601 = "2006-01-02"
)

type JSONResponse struct {
	Partners map[string][]PartnerUnprocessed `json:"partners"`
}

type PartnerUnprocessed struct {
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Email          string   `json:"email"`
	Country        string   `json:"country"`
	AvailableDates []string `json:"availableDates"`
}

type Partner struct {
	FirstName      string      `json:"firstName"`
	LastName       string      `json:"lastName"`
	Email          string      `json:"email"`
	Country        string      `json:"country"`
	AvailableDates []time.Time `json:"availableDates"`
}

type Country struct {
	AttendeeCount int      `json:"attendeeCount"`
	Attendees     []string `json:"attendees"`
	Name          string   `json:"name"`
	StartDate     string   `json:"startDate"`
}
