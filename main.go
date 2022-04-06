package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/google/go-github/v43/github"
	"github.com/rivo/tview"
)

// Create a struct representing a repository.
type repo struct {
	name        string
	description string
}

var (
	client *github.Client
	app    *tview.Application
	view   = View{}
)

func inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlQ:
		app.Stop()
	}
	return event
}

func main() {
	client = github.NewClient(nil)
	app = tview.NewApplication()

	rc := make(chan []*github.Repository)
	go fetchRepositories("akinsho", rc)
	go refreshRepositoryList(rc)
	grid := Layout()
	app.SetInputCapture(inputHandler)
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		log.Panicln(err)
	}
}
