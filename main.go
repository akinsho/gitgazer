package main

import (
	"akinsho/gitgazer/database"
	"akinsho/gitgazer/ui"
	"log"

	"akinsho/gitgazer/api"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := api.Setup()
	if err != nil {
		log.Panicln(err)
	}
	err = database.Setup()
	if err != nil {
		log.Panicln(err)
	}

	if err := ui.Setup(); err != nil {
		log.Panicln(err)
	}
}
