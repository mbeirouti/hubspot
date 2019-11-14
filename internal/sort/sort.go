package sort

import (
	"Hubspot/internal/models"
	"strings"
	"time"
)

type ByStartDate []string

func (timeSlot ByStartDate) Len() int      { return len(timeSlot) }
func (timeSlot ByStartDate) Swap(i, j int) { timeSlot[i], timeSlot[j] = timeSlot[j], timeSlot[i] }
func (timeSlot ByStartDate) Less(i, j int) bool {
	// Get two dates from each timeSlot string
	iTimeSlotDates := strings.Split(timeSlot[i], ",")
	jTimeSlotDates := strings.Split(timeSlot[j], ",")

	// Ignore errors because we assume that data is validated at this stage
	iStartDate, _ := time.Parse(models.ISO8601, iTimeSlotDates[0])
	jStartDate, _ := time.Parse(models.ISO8601, jTimeSlotDates[0])

	if iStartDate.Before(jStartDate) {
		return true
	}

	return false
}
