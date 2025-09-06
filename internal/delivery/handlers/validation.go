package handlers

import (
	"fmt"
	"regexp"
	"time"
)

var mmYYYYRe = regexp.MustCompile(`^(0[1-9]|1[0-2])-\d{4}$`)

func parseMonthYear(s string) (time.Time, error) {
	if !mmYYYYRe.MatchString(s) {
		return time.Time{}, fmt.Errorf("invalid format, expected MM-YYYY")
	}
	t, err := time.ParseInLocation("01-2006", s, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
