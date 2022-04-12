package main

import (
	"akinsho/gogazer/database"
	"akinsho/gogazer/ui"
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const tokenPath = "token.json"

func main() {
	token, err := retrieveAccessToken()
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.Token})
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
	err = database.Setup()
	if err != nil {
		log.Panicln(err)
	}

	if err := ui.Setup(client); err != nil {
		log.Panicln(err)
	}
}
