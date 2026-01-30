package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rzolkos/basecamp-cli/internal/config"
)

type InitCmd struct{}

func (c *InitCmd) Run(args []string) error {
	fmt.Fprintln(os.Stderr, "Basecamp CLI Configuration")
	fmt.Fprintln(os.Stderr, strings.Repeat("=", 40))

	reader := bufio.NewReader(os.Stdin)

	clientID := prompt(reader, "Client ID", "")
	clientSecret := prompt(reader, "Client Secret", "")
	accountID := prompt(reader, "Account ID", "")
	redirectURI := prompt(reader, "Redirect URI", "http://localhost:3002/callback")

	cfg := &config.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccountID:    accountID,
		RedirectURI:  redirectURI,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Configuration saved to: %s\n", config.ConfigFile())
	fmt.Fprintln(os.Stderr, "Run 'basecamp auth' to authenticate.")

	return PrintJSON(map[string]string{
		"status":  "ok",
		"message": "Configuration saved",
		"file":    config.ConfigFile(),
	})
}

func prompt(reader *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Fprintf(os.Stderr, "%s [%s]: ", label, defaultVal)
	} else {
		fmt.Fprintf(os.Stderr, "%s: ", label)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}
