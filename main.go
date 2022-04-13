package main

import (
	"akinsho/gitgazer/database"
	"akinsho/gitgazer/models"
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

	context := &models.GazeContext{
		Client: client,
	}

	if err := ui.Setup(context); err != nil {
		log.Panicln(err)
	}
}
