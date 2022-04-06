package main

import (
	"context"

	"github.com/google/go-github/v43/github"
)

var repositories []*github.Repository

func fetchRepositories(user string) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.List(
		context.Background(),
		user,
		&github.RepositoryListOptions{Sort: "updated"},
	)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// getRepositoryIssues fetches the issues for the given repository.
// using the github package to fetch the issues.
func getRepositoryIssues(repo *github.Repository) ([]*github.Issue, error) {
	issues, _, err := client.Issues.ListByRepo(
		context.Background(),
		repo.GetOwner().GetLogin(),
		repo.GetName(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func getRepositoryByIndex(index int) *github.Repository {
	if len(repositories) > 0 {
		return repositories[index]
	}
	return nil
}
