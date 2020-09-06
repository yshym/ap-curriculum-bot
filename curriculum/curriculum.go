// Package curriculum provides needed types and operations on curriculum data
package curriculum

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

// Week provides curriculum week data
type Week map[DayName]Day

// NewWeek creates new week object
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

type Meeting struct {
	Link string `json:"link"`
	Pass string `json:"pass"`
}
