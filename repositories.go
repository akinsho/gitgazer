package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v43/github"
)

var repositories []*github.Repository

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
	repositories = <-incomingRepos
	view.repoList.Clear()
	if len(repositories) == 0 {
		view.repoList.AddItem("No repositories found", "", 0, nil)
	}

	for _, repo := range repositories[:20] {
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

func getRepositoryByIndex(index int) *github.Repository {
	return repositories[index]
}
