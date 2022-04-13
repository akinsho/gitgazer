package main

import (
	gazerapp "akinsho/gitgazer/app"
	"akinsho/gitgazer/storage"
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
	db, err := storage.Setup()
	if err != nil {
		log.Panicln(err)
	}

	context := &gazerapp.Context{
		Client: client,
		DB:     db,
	}

	if err := ui.Setup(context); err != nil {
		log.Panicln(err)
	}
}
