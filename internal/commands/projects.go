package commands

import (
	"encoding/json"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

type ProjectsCmd struct{}

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
}

func (c *ProjectsCmd) Run(args []string) error {
	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/projects.json")
	if err != nil {
		return err
	}

	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return err
	}

	return PrintJSON(projects)
}
