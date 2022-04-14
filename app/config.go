package gazerapp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cli/oauth/api"
)

type Config struct {
	directory      string
	configFilepath string
	tokenPath      string
	Token          *api.AccessToken
}

const (
	configFile = "config.yaml"
	tokenFile  = "token.json"
	configDir  = "gitgazer"
)

// InitConfig setup the configuration file if need and read user options into state
// create the access token if required or read it from where it is stored
func InitConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	// TODO: investigate using https://github.com/adrg/xdg
	// to correctly derive the correct config paths for the user's OS
	dir := filepath.Join(home, ".config", configDir)
	config := &Config{
		directory:      dir,
		configFilepath: filepath.Join(dir, configFile),
		tokenPath:      filepath.Join(dir, tokenFile),
	}
	config.ensureDirectory()
	err = config.retrieveAccessToken()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) ensureDirectory() {
	if _, err := os.Stat(c.directory); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(c.directory, 0700)
	}
}

func (c *Config) persistToken(token *api.AccessToken) (err error) {
	jsonString, err := json.Marshal(token)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.tokenPath, jsonString, 0666); err != nil {
		return err
	}
	return nil
}

// readToken reads the token from the file and unmarshals it into a token struct
func (c *Config) readToken() (*api.AccessToken, error) {
	var token api.AccessToken
	jsonString, err := ioutil.ReadFile(c.tokenPath)
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
func (c *Config) retrieveAccessToken() error {
	var token *api.AccessToken
	if _, err := os.Stat(c.tokenPath); errors.Is(err, os.ErrNotExist) {
		token, err = getOAuthToken()
		if err != nil {
			return err
		}
		err = c.persistToken(token)
		if err != nil {
			return err
		}
	} else {
		token, err = c.readToken()
		if err != nil {
			return err
		}
	}
	c.Token = token
	return nil
}
