// Package curriculum provides needed types and operations on curriculum data
package curriculum

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

// Week provides curriculum week data
type Week map[DayName]Day

// NewWeek creates a Week object
func NewWeek() (*Week, error) {
	w := &Week{}

	filePath := path.Join("data", os.Getenv("CURRICULUM_FILE_NAME"))
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	err = w.FromJSON(f)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// FromJSON decodes json file into week object
func (w *Week) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(w)
}

// DayName is a name of the day of the week
type DayName string

// NewDayName creates a DayName object
func NewDayName(t *time.Time) DayName {
	return DayName(strings.ToLower(t.Weekday().String()))
}

// Day provides curriculum day data
type Day []DoublePeriodVariants

// DoublePeriodVariants provides all possible variants of current double period
type DoublePeriodVariants map[Date]DoublePeriod

// Date provides double period date
type Date string

// DoublePeriod provides curriculum day data
type DoublePeriod struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Lecturer string  `json:"lecturer"`
	Meeting  Meeting `json:"meeting"`
}

// Meeting provides meeting data
type Meeting struct {
	Link string `json:"link"`
	Pass string `json:"pass"`
}

// SpecificDay provides curriculum data of the specific day
type SpecificDay []DoublePeriod

// NewToday creates a Today object
func NewSpecificDay(w Week, t time.Time) SpecificDay {
	dayName := NewDayName(&t)

	day := w[dayName]

	date := Date(FormatTime(&t))

	var doublePeriods []DoublePeriod

	for _, dpv := range day {
		dp, _ := dpv[date]

		doublePeriods = append(doublePeriods, dp)
	}

	return SpecificDay(doublePeriods)
}

func Today(w Week) SpecificDay {
    return NewSpecificDay(w, time.Now())
}

// FormatTime returns formatted time
func FormatTime(t *time.Time) string {
	dayNumber, monthNumber := t.Day(), t.Month()
	var formattedTimeBuilder strings.Builder

	if dayNumber < 10 {
		formattedTimeBuilder.WriteString(fmt.Sprintf("0%d", dayNumber))
	} else {
		formattedTimeBuilder.WriteRune(rune(dayNumber))
	}

	formattedTimeBuilder.WriteRune('.')

	if monthNumber < 10 {
		formattedTimeBuilder.WriteString(fmt.Sprintf("0%d", monthNumber))
	} else {
		formattedTimeBuilder.WriteRune(rune(monthNumber))
	}

	return formattedTimeBuilder.String()
}
