package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yevhenshymotiuk/ap-curriculum-bot/curriculum"
)

func main() {
	w, err := curriculum.NewWeek()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(
		curriculum.NewSpecificDay(
			*w,
			time.Date(2020, 9, 7, 0, 0, 0, 0, time.Local),
		).Format(),
	)
}
