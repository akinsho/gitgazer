package main

import (
	"context"
	"fmt"
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

func fetchRepositories(user string, ch chan []*github.Repository) error {
	repos, _, err := client.Repositories.List(
		context.Background(),
		user,
		&github.RepositoryListOptions{Sort: "updated"},
	)
	if err != nil {
		fmt.Println(err)
	}
	ch <- repos
	return nil
}

func refreshRepositoryList(incomingRepos chan []*github.Repository) {
	repos := <-incomingRepos
	view.repoList.Clear()
	if len(repos) == 0 {
		view.repoList.AddItem("No repositories found", "", 0, nil)
	}

	for _, repo := range repos[:20] {
		name := repo.GetName()
		description := repo.GetDescription()
		if name != "" {
			showDesc := false
			if len(description) > 0 {
				showDesc = true
			}
			view.repoList.AddItem(repo.GetName(), description, 0, nil).ShowSecondaryText(showDesc)
		}
	}
	app.Draw()
}

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
