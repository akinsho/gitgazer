package main

import (
	"context"

	"github.com/google/go-github/v43/github"
)

var (
	repositories   []*github.Repository
	issuesByRepoID = make(map[int64][]*github.Issue)
)

func saveSelectedRepository(i int, s1, s2 string, r rune) (err error) {
	_, err = databaseConn.Insert(getRepositoryByIndex(i))
	return err
}

func fetchRepositories() ([]*github.Repository, error) {
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
	repoIssues := issuesByRepoID[repo.GetID()]
	if len(repoIssues) > 0 {
		return repoIssues, nil
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
	issuesByRepoID[repo.GetID()] = issues
	return issues, nil
}

func getRepositoryByIndex(index int) *github.Repository {
	if len(repositories) > 0 {
		return repositories[index]
	}
	return nil
}
