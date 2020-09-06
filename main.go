package main

import (
	"fmt"
	"log"

	"github.com/yevhenshymotiuk/ap-curriculum-bot/curriculum"
)

func main() {
	w, err := curriculum.NewWeek()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(w)
}
