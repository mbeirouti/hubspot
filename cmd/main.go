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

	// Put original data into a list of Partners instead (makes more sense to me)
	partners := []models.Partner{}
	for i, partner := range jsonResponse["partners"] {
		newPartner := models.Partner{
			FirstName:      partner.FirstName,
			LastName:       partner.LastName,
			Email:          partner.Email,
			Country:        partner.Country,
			AvailableDates: []time.Time{},
		}

		partners = append(partners, newPartner)
		for _, availability := range partner.AvailableDates {
			time, err := time.Parse(models.ISO8601, availability)
			if err != nil {
				fmt.Printf("error parsing time (%s) %s", availability, err)
				return
			}

			partners[i].AvailableDates = append(partners[i].AvailableDates, time)
		}
	}

	// Group partners by country in a new map
	partnersByCountry := map[string][]models.Partner{}
	for _, partner := range partners {
		if _, ok := partnersByCountry[partner.Country]; ok {
			partnersByCountry[partner.Country] = append(partnersByCountry[partner.Country], partner)
		} else {
			partnersByCountry[partner.Country] = []models.Partner{}
			partnersByCountry[partner.Country] = append(partnersByCountry[partner.Country], partner)
		}
	}

	timeSlotsByCountry := map[string]map[string]int{}
	partnersByTimeSlot := map[string]map[string][]string{}
	for country := range partnersByCountry {
		partnersByTimeSlot[country] = map[string][]string{}

		timeSlotCount := map[string]int{}
		// For each partner in each country
		for _, partner := range partnersByCountry[country] {
			// For all available dates of the partner
			for i := 0; i < len(partner.AvailableDates)-1; i++ {
				// Find all the two consecutive day slots (assuming dates are already sorted)
				timeDifference := partner.AvailableDates[i+1].Sub(partner.AvailableDates[i]).Hours()
				if int(timeDifference) == oneDay {
					// Add partner to the count of partners with that two day slot available
					// Each available two day slot is represented by the unique string slotString
					slotString := fmt.Sprintf("%s,%s", partner.AvailableDates[i].Format(models.ISO8601), partner.AvailableDates[i+1].Format(models.ISO8601))
					if _, ok := timeSlotCount[slotString]; ok {
						timeSlotCount[slotString] += 1
					} else {
						timeSlotCount[slotString] = 1
					}

					// Add partner to the list of partners available for that two day slot
					if _, ok := partnersByTimeSlot[country][slotString]; ok {
						partnersByTimeSlot[country][slotString] = append(partnersByTimeSlot[country][slotString], partner.Email)
					} else {
						partnersByTimeSlot[country][slotString] = []string{}
						partnersByTimeSlot[country][slotString] = append(partnersByTimeSlot[country][slotString], partner.Email)
					}
				}
			}
		}

		timeSlotsByCountry[country] = timeSlotCount
	}

	countryMaxTimeSlot := map[string]string{}
	// Find max time slot by country
	for country := range timeSlotsByCountry {
		maxCount := -1
		for _, count := range timeSlotsByCountry[country] {
			if count > maxCount {
				maxCount = count
			}
		}

		// Get all time slots with the maxCount
		maxTimeSlots := []string{}
		for timeSlot, count := range timeSlotsByCountry[country] {
			if count == maxCount {
				maxTimeSlots = append(maxTimeSlots, timeSlot)
			}
		}

		// Get the earliest time slot with maxCount
		sort.Sort(mysort.ByTime(maxTimeSlots))

		// Set that as the correct answer
		countryMaxTimeSlot[country] = maxTimeSlots[0]
	}

	// Initialize response
	response := map[string][]models.Country{}
	response["countries"] = []models.Country{}

	// Add necessary data to response
	for country, maxTimeSlot := range countryMaxTimeSlot {
		// Get start date only from the timeslot
		dates := strings.Split(maxTimeSlot, ",")
		startDate := dates[0]

		// Append partners that are available at maxTimeSlot
		newCountry := models.Country{
			AttendeeCount: len(partnersByTimeSlot[country][maxTimeSlot]),
			Attendees:     partnersByTimeSlot[country][maxTimeSlot],
			Name:          country,
			StartDate:     startDate,
		}

		response["countries"] = append(response["countries"], newCountry)
	}

	// Post Response
	resp, err := requests.PostData(responseURL, response)
	if err != nil {
		fmt.Printf("error posting response %s", err)
		return
	}

	// Check response
	if resp.StatusCode != http.StatusOK {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body %s", err)
			return
		}

		fmt.Println(resp.StatusCode)
		fmt.Println(string(respBytes))
	}
}
