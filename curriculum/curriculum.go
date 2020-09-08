// Package curriculum provides needed types and operations on curriculum data
package curriculum

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/yevhenshymotiuk/ap-curriculum-bot/helpers"
)

// Week provides curriculum week data
type Week map[helpers.DayName]Day

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

// Day provides curriculum day data
type Day []SubGroup

type SubGroup [2]DoublePeriodVariants

// DoublePeriodVariants provides all possible variants of current double period
type DoublePeriodVariants map[helpers.Date]DoublePeriod

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
type SpecificDay [2][]DoublePeriod

// NewSpecificDay creates a SpecificDay object
func NewSpecificDay(w Week, t time.Time) SpecificDay {
	dayName := helpers.NewDayName(&t)

	day := w[dayName]

	date := helpers.Date(helpers.FormatTime(&t))

	var doublePeriods1, doublePeriods2 []DoublePeriod

	for _, sg := range day {
		var dpv1, dpv2 DoublePeriodVariants

		if len(sg[1]) == 0 {
			dpv1, dpv2 = sg[0], sg[0]
		} else {
			dpv1, dpv2 = sg[0], sg[1]
		}

		dp1, _ := dpv1[date]
		dp2, _ := dpv2[date]

		doublePeriods1 = append(doublePeriods1, dp1)
		doublePeriods2 = append(doublePeriods2, dp2)
	}

	return SpecificDay([2][]DoublePeriod{doublePeriods1, doublePeriods2})
}

func (sd SpecificDay) Format() string {
	var formatted string

	dps1, dps2 := sd[0], sd[1]

	var fdps1, fdps2 string
	fdps1 = formatDoublePeriods(dps1)

	if reflect.DeepEqual(dps1, dps2) {
		formatted = fmt.Sprintf(
			"Розклад однаковий для обох підгруп:\n%s",
			fdps1,
		)
	} else {
		fdps2 = formatDoublePeriods(dps2)
		formatted = fmt.Sprintf("Підгрупа 1:\n%s\n\nПідгрупа 2:\n%s", fdps1, fdps2)
	}

	return formatted
}

func formatDoublePeriods(dps []DoublePeriod) string {
	var formatted strings.Builder

	if len(dps) == 0 {
		return "Пар немає"
	}

	for i, dp := range dps {
		if dp == (DoublePeriod{}) {
			formatted.WriteString(fmt.Sprintf("%d) -\n", i+1))
		} else {
			formatted.WriteString(
				fmt.Sprintf(
					"%d) %s(%s) | %s\nПосилання: %s\n",
					i+1,
					dp.Name,
					dp.Type,
					dp.Lecturer,
					dp.Meeting.Link,
				),
			)

			pass := dp.Meeting.Pass
			if pass != "" {
				formatted.WriteString(fmt.Sprintf("Пароль: %s\n", pass))
			}
		}
	}

	return strings.TrimSpace(formatted.String())
}

// Today creates a SpecificDay object for today's curriculum
func Today(w Week) SpecificDay {
	l, err := helpers.LoadLocation()
	if err != nil {
		log.Fatalln(err)
	}

	return NewSpecificDay(w, time.Now().In(l))
}
