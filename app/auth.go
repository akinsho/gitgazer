package app

import (
	"os"

	"github.com/cli/oauth"
	"github.com/cli/oauth/api"
)

// getOAuthToken authenticate the user with Github and return an access token
func getOAuthToken() (*api.AccessToken, error) {
	flow := &oauth.Flow{
		Host:         oauth.GitHubHost("https://github.com"),
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		CallbackURI:  "http://127.0.0.1/callback",
		Scopes:       []string{"repo", "read:org", "gist"},
	}

	accessToken, err := flow.DetectFlow()
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}
