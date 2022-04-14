package api

import (
	"akinsho/gitgazer/domain"
	"context"

	"github.com/cli/oauth/api"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client struct {
	graphql *githubv4.Client
}

func Setup(token *api.AccessToken) (*Client, error) {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.Token})
	httpClient := oauth2.NewClient(context.Background(), src)
	return &Client{githubv4.NewClient(httpClient)}, nil
}

func (c *Client) ListStarredRepositories() ([]*domain.Repository, error) {
	var starredRepositoriesQuery struct {
		Viewer struct {
			StarredRepositories struct {
				Nodes []*domain.Repository `graphql:"nodes"`
			} `graphql:"starredRepositories(first: $repoCount, orderBy: {field: STARRED_AT, direction: DESC})"`
		}
	}

	err := c.graphql.Query(
		context.Background(),
		&starredRepositoriesQuery,
		map[string]interface{}{
			"labelCount": githubv4.Int(20),
			"issueCount": githubv4.Int(20),
			"repoCount":  githubv4.Int(20),
			"issuesOrderBy": githubv4.IssueOrder{
				Direction: githubv4.OrderDirectionDesc,
				Field:     githubv4.IssueOrderFieldUpdatedAt,
			},
			"prCount": githubv4.Int(5),
			"prState": []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
			"pullRequestOrderBy": githubv4.IssueOrder{
				Direction: githubv4.OrderDirectionDesc,
				Field:     githubv4.IssueOrderFieldUpdatedAt,
			},
		},
	)
	return starredRepositoriesQuery.Viewer.StarredRepositories.Nodes, err
}

func (c *Client) FetchRepositoryByName(name, owner string) (*domain.Repository, error) {

	var repositoryQuery struct {
		Repository domain.Repository `graphql:"repository(name: $name, owner: $owner)"`
	}
	variables := map[string]interface{}{
		"name":       githubv4.String(name),
		"owner":      githubv4.String(owner),
		"labelCount": githubv4.Int(20),
		"issueCount": githubv4.Int(20),
		"issuesOrderBy": githubv4.IssueOrder{
			Direction: githubv4.OrderDirectionDesc,
			Field:     githubv4.IssueOrderFieldUpdatedAt,
		},
		"prCount": githubv4.Int(5),
		"prState": []githubv4.PullRequestState{githubv4.PullRequestStateOpen},
		"pullRequestOrderBy": githubv4.IssueOrder{
			Direction: githubv4.OrderDirectionDesc,
			Field:     githubv4.IssueOrderFieldUpdatedAt,
		},
	}
	err := c.graphql.Query(context.Background(), &repositoryQuery, variables)
	return &repositoryQuery.Repository, err
}
