package main

import (
	gazerapp "akinsho/gitgazer/app"
	"akinsho/gitgazer/database"
	"akinsho/gitgazer/ui"
	"log"

	"akinsho/gitgazer/api"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	client, err := api.Setup()
	if err != nil {
		log.Panicln(err)
	}
	err = database.Setup()
	if err != nil {
		log.Panicln(err)
	}

	context := &gazerapp.Context{
		Client: client,
	}

	if err := ui.Setup(context); err != nil {
		log.Panicln(err)
	}
}
