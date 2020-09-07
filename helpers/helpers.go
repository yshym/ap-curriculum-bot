package helpers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func LoadLocation() (*time.Location, error) {
	return time.LoadLocation("Europe/Kiev")
}

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
		formattedTimeBuilder.WriteString(monthNumber.String())
	}

	return formattedTimeBuilder.String()
}
