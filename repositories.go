package main

import (
	"context"

	"github.com/google/go-github/v43/github"
)

var (
	repositories []*github.Repository
	issues       []*github.Issue
)

func fetchRepositories(user string) ([]*github.Repository, error) {
	//  TODO: We need a way to invalidate previous fetched repositories
	// and refetch but this is necessary for now to prevent DDOSing the API.
	if len(repositories) > 0 {
		return repositories, nil
	}

	repos, _, err := client.Repositories.List(
		context.Background(),
		user,
		&github.RepositoryListOptions{Sort: "updated"},
	)
	if err != nil {
		return nil, err
	}
	repositories = repos
	return repos, nil
}

// getRepositoryIssues fetches the issues for the given repository.
// using the github package to fetch the issues.
func getRepositoryIssues(repo *github.Repository) ([]*github.Issue, error) {
	//  TODO: Cache invalidation
	if len(issues) > 0 {
		return issues, nil
	}
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
