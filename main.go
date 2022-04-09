package main

import (
	"akinsho/gogazer/database"
	"context"
	"log"
	"os"

	"github.com/cli/oauth"
	"github.com/cli/oauth/api"
	"github.com/gdamore/tcell/v2"
	"github.com/google/go-github/v43/github"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"
)

// Create a struct representing a repository.
type repo struct {
	name        string
	description string
}

var (
	client       *github.Client
	app          *tview.Application
	databaseConn *database.Gazers
	view         = View{}
)

func inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		app.Stop()
	}
	return event
}

// getOAuthToken authenticate the user with Github and return an access token
func getOAuthToken() (*api.AccessToken, error) {
	flow := &oauth.Flow{
		Host:         oauth.GitHubHost("https://github.com"),
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		CallbackURI:  "http://127.0.0.1/callback",
		Scopes:       []string{"repo", "read:org", "gist"},
	}

	accessToken, err := flow.DetectFlow()
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func main() {
	token, err := getOAuthToken()
	if err != nil {
		log.Panicln(err)
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.Token})
	httpClient := oauth2.NewClient(context.Background(), src)
	client = github.NewClient(httpClient)
	app = tview.NewApplication()
	db, err := database.Setup()
	if err != nil {
		log.Panicln(err)
	}
	databaseConn = db

	refreshRepositoryList()
	layout := getLayout()
	layout.SetTitle("Go gazer")
	app.SetInputCapture(inputHandler)
	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		log.Panicln(err)
	}
}
