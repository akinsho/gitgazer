package main

import (
	"akinsho/gitgazer/app"
	"akinsho/gitgazer/storage"
	"akinsho/gitgazer/ui"
	"log"

	"akinsho/gitgazer/api"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config, err := app.InitConfig()
	if err != nil {
		log.Panicln(err)
	}
	client, err := api.Setup(config.Token)
	if err != nil {
		log.Panicln(err)
	}
	db, err := storage.Setup()
	if err != nil {
		log.Panicln(err)
	}

	context := &app.Context{
		Client: client,
		Config: config,
		DB:     db,
	}

	if err := ui.Setup(context); err != nil {
		log.Panicln(err)
	}
}
