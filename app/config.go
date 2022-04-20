package app

import (
	"akinsho/gitgazer/domain"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cli/oauth/api"
	"gopkg.in/yaml.v2"
)

type LogOptions struct {
	Enabled bool `yaml:"enabled"`
}

type PanelDetails struct {
	Preferred domain.PanelName `yaml:"preferred"`
}

type Panels struct {
	Details PanelDetails `yaml:"details"`
	Log     LogOptions   `yaml:"log"`
}

type UserConfig struct {
	Panels Panels `yaml:"panels"`
}

type Config struct {
	directory      string
	configFilepath string
	tokenPath      string
	StoragePath    string
	Token          *api.AccessToken
	UserConfig     *UserConfig
}

const (
	configFile  = "config.yaml"
	tokenFile   = "token.json"
	configDir   = "gitgazer"
	StoragePath = "gazers.db"
)

var defaults = &Config{
	UserConfig: &UserConfig{
		Panels: Panels{
			Log: LogOptions{
				Enabled: false,
			},
			Details: PanelDetails{
				Preferred: domain.PullRequestPanel,
			},
		},
	},
}

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
		StoragePath:    filepath.Join(dir, StoragePath),
	}
	err = config.ensureDirectory()
	if err != nil {
		return nil, err
	}
	if !config.exists() {
		config.UserConfig, err = writeConfig(config.configFilepath, defaults.UserConfig)
	} else {
		config.UserConfig, err = readConfig(config.configFilepath, defaults.UserConfig)
	}
	if err != nil {
		return nil, err
	}
	err = config.retrieveAccessToken()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// writeConfig writes the default config file to the config directory
func writeConfig(path string, def *UserConfig) (*UserConfig, error) {
	file, err := os.Create(path)
	if err != nil {
		return def, err
	}
	e := yaml.NewEncoder(file)
	if err := e.Encode(def); err != nil {
		return def, err
	}
	defer file.Close()
	return def, nil
}

// readConfig returns a new decoded Config struct
func readConfig(path string, config *UserConfig) (*UserConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) ensureDirectory() error {
	if _, err := os.Stat(c.directory); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(c.directory, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) exists() bool {
	if _, err := os.Stat(c.configFilepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
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
