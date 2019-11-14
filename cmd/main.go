package main

import (
	"Hubspot/internal/models"
	"Hubspot/internal/requests"
	mysort "Hubspot/internal/sort"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	oneDay = 24
)

var (
	apiURL      = "https://candidate.hubteam.com/candidateTest/v3/problem/dataset?userKey=23b677b106e1bb467840d4659d5d"
	responseURL = "https://candidate.hubteam.com/candidateTest/v3/problem/result?userKey=23b677b106e1bb467840d4659d5d"
)

func main() {

	// Get JSON Data and decode into partners map which only contains one key "partners"
	jsonResponse := map[string][]models.PartnerUnprocessed{}
	err := requests.GetData(apiURL, &jsonResponse)
	if err != nil {
		fmt.Printf("error building request %s", err)
	}

	// Put partners into a dictionary mapping country to list of partners for that country
	partnersByCountry := map[string][]models.Partner{}
	for _, partner := range jsonResponse["partners"] {
		newPartner := models.Partner{
			FirstName:      partner.FirstName,
			LastName:       partner.LastName,
			Email:          partner.Email,
			Country:        partner.Country,
			AvailableDates: []time.Time{},
		}

		// Parse date strings into actual Golang time.Time objects
		// This is a workaround to an issue inherent to Golang's JSON decoding library
		for _, availability := range partner.AvailableDates {
			time, err := time.Parse(models.ISO8601, availability)
			if err != nil {
				fmt.Printf("error parsing time (%s) %s", availability, err)
				return
			}

			newPartner.AvailableDates = append(newPartner.AvailableDates, time)
		}

		partnersByCountry[partner.Country] = append(partnersByCountry[partner.Country], newPartner)
	}

	// For each country make map with each valid timeSlot and a list of the partners available at that timeSlot
	partnersByTimeSlot := map[string]map[string][]string{}
	for country := range partnersByCountry {
		partnersByTimeSlot[country] = map[string][]string{}
		// For each partner in each country
		for _, partner := range partnersByCountry[country] {
			// For all available dates of the partner
			for i := 0; i < len(partner.AvailableDates)-1; i++ {
				// Find all the two consecutive day slots (assuming dates are already sorted)
				timeDifference := partner.AvailableDates[i+1].Sub(partner.AvailableDates[i]).Hours()
				if int(timeDifference) == oneDay {
					// Add partner to the list of partners available for that two day slot
					// Each available two day slot is represented by a unique string called timeSlot
					timeSlot := fmt.Sprintf("%s,%s", partner.AvailableDates[i].Format(models.ISO8601), partner.AvailableDates[i+1].Format(models.ISO8601))
					partnersByTimeSlot[country][timeSlot] = append(partnersByTimeSlot[country][timeSlot], partner.Email)
				}
			}
		}
	}

	// Find  earliest two day timeSlot with max number of attending partners per country
	countryMaxTimeSlot := map[string]string{}
	for country := range partnersByTimeSlot {

		maxCount := -1
		for _, partners := range partnersByTimeSlot[country] {
			if len(partners) > maxCount {
				maxCount = len(partners)
			}
		}

		// Get all timeSlots with the maxCount (this could be made more efficient but at the cost of simplicity)
		maxTimeSlots := []string{}
		for timeSlot, partners := range partnersByTimeSlot[country] {
			if len(partners) == maxCount {
				maxTimeSlots = append(maxTimeSlots, timeSlot)
			}
		}

		// Get the earliest start date timeSlot with maxCount
		sort.Sort(mysort.ByStartDate(maxTimeSlots))

		// Set that as the correct answer
		countryMaxTimeSlot[country] = maxTimeSlots[0]
	}

	// Initialize result
	response := map[string][]models.Country{}
	response["countries"] = []models.Country{}

	// Add necessary data to result
	for country, maxTimeSlot := range countryMaxTimeSlot {
		// Get start date only from the timeslot
		dates := strings.Split(maxTimeSlot, ",")
		startDate := dates[0]

		// Add all partners available at maxTimeSlot
		newCountry := models.Country{
			AttendeeCount: len(partnersByTimeSlot[country][maxTimeSlot]),
			Attendees:     partnersByTimeSlot[country][maxTimeSlot],
			Name:          country,
			StartDate:     startDate,
		}

		response["countries"] = append(response["countries"], newCountry)
	}

	// Post result
	resp, err := requests.PostData(responseURL, response)
	if err != nil {
		fmt.Printf("error posting response %s", err)
		return
	}

	// Check response for posted result
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("You did it.. Woot!")
	case http.StatusBadRequest:
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body %s", err)
			return
		}

		fmt.Println(resp.StatusCode)
		fmt.Println(string(respBytes))
	default:
		fmt.Println("Uh oh, spadoodios!")
	}
}
