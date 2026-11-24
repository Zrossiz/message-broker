package main

import (
	"log"

	"github.com/Zrossiz/message-broker/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Stop()
	if err != nil {
		log.Fatal(err)
	}
}
