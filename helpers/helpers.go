package helpers

import (
	"log"
	"time"
)

func LoadLocation() (*time.Location, error) {
	return time.LoadLocation("Europe/Kiev")
}

func Now() (time.Time) {
	l, err := LoadLocation()
	if err != nil {
		log.Fatalln(err)
	}

    return time.Now().In(l)
}
