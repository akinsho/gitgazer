package api

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client struct {
	Graphql *githubv4.Client
}

func Setup() (*Client, error) {
	token, err := retrieveAccessToken()
	if err != nil {
		return nil, err
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.Token})
	httpClient := oauth2.NewClient(context.Background(), src)
	return &Client{githubv4.NewClient(httpClient)}, nil
}
