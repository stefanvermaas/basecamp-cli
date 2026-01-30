package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConfigSaveLoad(t *testing.T) {
	// Use temp directory for test
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	cfg := &Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		AccountID:    "12345",
		RedirectURI:  "http://localhost:3002/callback",
	}

	// Save config
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tmpDir, "basecamp", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// Load config
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.ClientID != cfg.ClientID {
		t.Errorf("ClientID = %v, want %v", loaded.ClientID, cfg.ClientID)
	}
	if loaded.ClientSecret != cfg.ClientSecret {
		t.Errorf("ClientSecret = %v, want %v", loaded.ClientSecret, cfg.ClientSecret)
	}
	if loaded.AccountID != cfg.AccountID {
		t.Errorf("AccountID = %v, want %v", loaded.AccountID, cfg.AccountID)
	}
}

func TestConfigNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	_, err := Load()
	if err != ErrConfigNotFound {
		t.Errorf("Load() error = %v, want %v", err, ErrConfigNotFound)
	}
}

func TestAPIBaseURL(t *testing.T) {
	cfg := &Config{AccountID: "12345"}
	want := "https://3.basecampapi.com/12345"
	if got := cfg.APIBaseURL(); got != want {
		t.Errorf("APIBaseURL() = %v, want %v", got, want)
	}
}

func TestGetRedirectURI(t *testing.T) {
	tests := []struct {
		name        string
		redirectURI string
		want        string
	}{
		{"default", "", "http://localhost:3002/callback"},
		{"custom", "http://localhost:4000/cb", "http://localhost:4000/cb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{RedirectURI: tt.redirectURI}
			if got := cfg.GetRedirectURI(); got != tt.want {
				t.Errorf("GetRedirectURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_DATA_HOME", tmpDir)
	defer os.Unsetenv("XDG_DATA_HOME")

	token := &TokenData{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresIn:    3600,
	}

	// Save token
	if err := SaveToken(token); err != nil {
		t.Fatalf("SaveToken() error = %v", err)
	}

	// Verify file exists
	tokenPath := filepath.Join(tmpDir, "basecamp", "token.json")
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		t.Fatal("token file was not created")
	}

	// Load token
	accessToken, err := LoadToken()
	if err != nil {
		t.Fatalf("LoadToken() error = %v", err)
	}

	if accessToken != token.AccessToken {
		t.Errorf("LoadToken() = %v, want %v", accessToken, token.AccessToken)
	}
}

func TestTokenNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_DATA_HOME", tmpDir)
	defer os.Unsetenv("XDG_DATA_HOME")

	_, err := LoadToken()
	if err != ErrNotAuthenticated {
		t.Errorf("LoadToken() error = %v, want %v", err, ErrNotAuthenticated)
	}
}

func TestTokenExpired(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_DATA_HOME", tmpDir)
	defer os.Unsetenv("XDG_DATA_HOME")

	token := &TokenData{
		AccessToken: "expired-token",
		ExpiresAt:   time.Now().Unix() - 3600, // expired 1 hour ago
	}

	tokenDir := filepath.Join(tmpDir, "basecamp")
	os.MkdirAll(tokenDir, 0700)
	data, _ := json.MarshalIndent(token, "", "  ")
	tokenPath := filepath.Join(tokenDir, "token.json")
	os.WriteFile(tokenPath, data, 0600)

	_, err := LoadToken()
	if err != ErrTokenExpired {
		t.Errorf("LoadToken() error = %v, want %v", err, ErrTokenExpired)
	}
}
