package main

import (
	"log"
	"os"

	"github.com/Zrossiz/message-broker/internal/app"
)

func main() {
	err := app.New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
