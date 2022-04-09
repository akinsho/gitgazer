package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

func persistToken(token *api.AccessToken) (err error) {
	jsonString, err := json.Marshal(token)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tokenPath, jsonString, 0666); err != nil {
		return err
	}
	return nil
}

// readToken reads the token from the file and unmarshals it into a token struct
func readToken() (*api.AccessToken, error) {
	var token api.AccessToken
	jsonString, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonString, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// retrieveAccessToken if an access token has been saved previously then read it back
// into memory from the file where it is saved otherwise start a new oauth flow and persist
// the token to the file
func retrieveAccessToken() (*api.AccessToken, error) {
	var token *api.AccessToken
	if _, err := os.Stat(tokenPath); errors.Is(err, os.ErrNotExist) {
		token, err = getOAuthToken()
		if err != nil {
			return nil, err
		}
		err = persistToken(token)
		if err != nil {
			return nil, err
		}
	} else {
		token, err = readToken()
		if err != nil {
			return nil, err
		}
	}
	return token, nil
}
