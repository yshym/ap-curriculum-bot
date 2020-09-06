// Package curriculum provides needed types and operations on curriculum data
package curriculum

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

// Week provides curriculum week data
type Week map[DayName]Day

// NewWeek creates a Week object
func NewWeek(r io.Reader) (*Week, error) {
	w := &Week{}

	err := w.FromJSON(r)
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

func (sd SpecificDay) Format() string {
	var formatted strings.Builder

	for i, dp := range sd {
		if dp != (DoublePeriod{}) {
			formatted.WriteString(
				fmt.Sprintf(
					"%d) %s(%s) | %s\n%s\n",
					i+1,
					dp.Name,
					dp.Type,
					dp.Lecturer,
					dp.Meeting.Link,
				),
			)
		}
	}

	return strings.TrimSpace(formatted.String())
}

func Today(w Week) SpecificDay {
	l, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		log.Fatalln(err)
	}

	return NewSpecificDay(w, time.Now().In(l))
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
