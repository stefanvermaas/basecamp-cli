package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/rzolkos/basecamp-cli/internal/config"
)

type AuthCmd struct{}

func (c *AuthCmd) Run(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Basecamp OAuth Authentication")
	fmt.Fprintln(os.Stderr, "========================================")

	// Parse redirect URI to get port
	redirectURL, err := url.Parse(cfg.GetRedirectURI())
	if err != nil {
		return fmt.Errorf("invalid redirect URI: %w", err)
	}
	port := redirectURL.Port()
	if port == "" {
		port = "3002"
	}

	// Channel to receive auth code
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Start callback server
	server := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code != "" {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `<html><body style="font-family:sans-serif;text-align:center;padding:50px;">
<h1>Authentication Successful!</h1><p>You can close this window.</p></body></html>`)
				codeCh <- code
			} else {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `<html><body style="font-family:sans-serif;text-align:center;padding:50px;">
<h1>Authentication Failed</h1><p>No authorization code received.</p></body></html>`)
				errCh <- fmt.Errorf("no authorization code received")
			}
		}),
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start callback server: %w", err)
	}

	fmt.Fprintf(os.Stderr, "\nStarting callback server on port %s...\n", port)

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Build and open authorization URL
	params := url.Values{
		"type":         {"web_server"},
		"client_id":    {cfg.ClientID},
		"redirect_uri": {cfg.GetRedirectURI()},
	}
	authURL := config.AuthorizationURL + "?" + params.Encode()

	fmt.Fprintln(os.Stderr, "\nOpening browser for authorization...")
	fmt.Fprintf(os.Stderr, "URL: %s\n", authURL)
	fmt.Fprintln(os.Stderr, "\nIf browser doesn't open, copy the URL above.")

	openBrowser(authURL)

	// Wait for callback with timeout
	fmt.Fprintln(os.Stderr, "\nWaiting for authorization...")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var authCode string
	select {
	case authCode = <-codeCh:
		fmt.Fprintln(os.Stderr, "Authorization code received")
	case err := <-errCh:
		server.Shutdown(context.Background())
		return err
	case <-ctx.Done():
		server.Shutdown(context.Background())
		return fmt.Errorf("timeout waiting for authorization")
	}

	server.Shutdown(context.Background())

	// Exchange code for token
	fmt.Fprintln(os.Stderr, "\nExchanging code for token...")

	token, err := exchangeCodeForToken(cfg, authCode)
	if err != nil {
		return err
	}

	if err := config.SaveToken(token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	fmt.Fprintln(os.Stderr, "\nAuthentication successful!")
	fmt.Fprintf(os.Stderr, "Token saved to: %s\n", config.TokenFile())

	return PrintJSON(map[string]string{
		"status":  "ok",
		"message": "Authentication successful",
		"file":    config.TokenFile(),
	})
}

func exchangeCodeForToken(cfg *config.Config, code string) (*config.TokenData, error) {
	data := url.Values{
		"type":          {"web_server"},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"redirect_uri":  {cfg.GetRedirectURI()},
		"code":          {code},
	}

	resp, err := http.PostForm(config.TokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %d %s\n%s", resp.StatusCode, resp.Status, string(body))
	}

	var token config.TokenData
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return &token, nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return
	}

	cmd.Start()
}
