package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	AuthorizationURL = "https://launchpad.37signals.com/authorization/new"
	TokenURL         = "https://launchpad.37signals.com/authorization/token"
)

var (
	ErrConfigNotFound   = errors.New("config file not found, run 'basecamp init' to create one")
	ErrNotAuthenticated = errors.New("not authenticated, run 'basecamp auth' first")
	ErrTokenExpired     = errors.New("token expired, run 'basecamp auth' to refresh")
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccountID    string `json:"account_id"`
	RedirectURI  string `json:"redirect_uri"`
}

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
}

func configDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "basecamp")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "basecamp")
}

func dataDir() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, "basecamp")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "basecamp")
}

func ConfigFile() string {
	return filepath.Join(configDir(), "config.json")
}

func TokenFile() string {
	return filepath.Join(dataDir(), "token.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(ConfigFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path := ConfigFile()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (c *Config) APIBaseURL() string {
	return "https://3.basecampapi.com/" + c.AccountID
}

func (c *Config) GetRedirectURI() string {
	if c.RedirectURI == "" {
		return "http://localhost:3002/callback"
	}
	return c.RedirectURI
}

func LoadToken() (string, error) {
	data, err := os.ReadFile(TokenFile())
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotAuthenticated
		}
		return "", err
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return "", err
	}

	if token.ExpiresAt > 0 && time.Now().Unix() > token.ExpiresAt {
		return "", ErrTokenExpired
	}

	return token.AccessToken, nil
}

func SaveToken(token *TokenData) error {
	if token.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Unix() + token.ExpiresIn
	}

	path := TokenFile()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
