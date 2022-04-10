package main

import (
	"akinsho/gogazer/database"
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const tokenPath = "token.json"

var (
	client       *githubv4.Client
	app          *tview.Application
	databaseConn *database.Gazers
	view         = View{}
)

func main() {
	token, err := retrieveAccessToken()
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.Token})
	httpClient := oauth2.NewClient(context.Background(), src)
	client = githubv4.NewClient(httpClient)
	app = tview.NewApplication()
	db, err := database.Setup()
	if err != nil {
		log.Panicln(err)
	}
	databaseConn = db

	go refreshRepositoryList()
	layout := getLayout()
	layout.SetTitle("Go gazer")
	app.SetInputCapture(appInputHandler)
	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		log.Panicln(err)
	}
}
