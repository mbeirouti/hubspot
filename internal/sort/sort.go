package sort

import (
	"Hubspot/internal/models"
	"strings"
	"time"
)

type ByTime []string

func (t ByTime) Len() int      { return len(t) }
func (t ByTime) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ByTime) Less(i, j int) bool {
	tis := strings.Split(t[i], ",")
	tjs := strings.Split(t[j], ",")

	tiTime, _ := time.Parse(models.ISO8601, tis[0])
	tjTime, _ := time.Parse(models.ISO8601, tjs[0])

	if tiTime.Before(tjTime) {
		return true
	}

	return false
}
