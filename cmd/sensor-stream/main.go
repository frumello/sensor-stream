package main

import (
	"log"
	"sensor-stream/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal("error configuring application")
	}
	if err = a.Process(); err != nil {
		log.Fatal("error running application", err)
	}
}
