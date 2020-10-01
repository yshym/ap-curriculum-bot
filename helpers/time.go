package helpers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// LoadLocation loads the specific location
func LoadLocation() (*time.Location, error) {
	return time.LoadLocation("Europe/Kiev")
}

// Now returns time.Time objects relying on the specific location
func Now() time.Time {
	l, err := LoadLocation()
	if err != nil {
		log.Fatalln(err)
	}

	return time.Now().In(l)
}

// FormatTime returns formatted time
func FormatTime(t *time.Time) string {
	dayNumber, monthNumber := t.Day(), t.Month()
	var formattedTimeBuilder strings.Builder

	if dayNumber < 10 {
		formattedTimeBuilder.WriteString(fmt.Sprintf("0%d", dayNumber))
	} else {
		formattedTimeBuilder.WriteString(strconv.Itoa(dayNumber))
	}

	formattedTimeBuilder.WriteRune('.')

	if monthNumber < 10 {
		formattedTimeBuilder.WriteString(fmt.Sprintf("0%d", monthNumber))
	} else {
		formattedTimeBuilder.WriteString(strconv.Itoa(int(monthNumber)))
	}

	return formattedTimeBuilder.String()
}

// DayName is a name of the day of the week
type DayName string

// NewDayName creates a DayName object
func NewDayName(t *time.Time) DayName {
	return DayName(strings.ToLower(t.Weekday().String()))
}

// Date provides double period date
type Date string
