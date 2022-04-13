package api

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var Client *githubv4.Client

func Setup() error {
	token, err := retrieveAccessToken()
	if err != nil {
		return err
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.Token})
	httpClient := oauth2.NewClient(context.Background(), src)
	Client = githubv4.NewClient(httpClient)
	return nil
}
