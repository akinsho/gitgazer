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

func fetchRepositories(user string, channel chan []*github.Repository) error {
	repos, _, err := client.Repositories.List(
		context.Background(),
		user,
		&github.RepositoryListOptions{Sort: "updated"},
	)
	if err != nil {
		fmt.Println(err)
	}
	channel <- repos
	return nil
}

func refreshRepositoryList(incomingRepos chan []*github.Repository) {
	repos := <-incomingRepos
	view.repoList.Clear()
	if len(repos) == 0 {
		view.repoList.AddItem("No repositories found", "", 0, nil)
	}

	for _, repo := range repos[:20] {
		if repo.Name != nil {
			view.repoList.AddItem(repo.GetName(), repo.GetDescription(), 0, nil)
		}
	}
	app.Draw()
}

func main() {
	client = github.NewClient(nil)
	app = tview.NewApplication()

	reposChan := make(chan []*github.Repository)
	go fetchRepositories("akinsho", reposChan)
	grid := Layout()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.Stop()
		}
		return event
	})
	go refreshRepositoryList(reposChan)
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		log.Panicln(err)
	}
}
