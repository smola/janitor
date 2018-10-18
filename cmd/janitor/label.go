package main

import (
	"context"

	"gopkg.in/src-d/go-cli.v0"

	"github.com/smola/janitor/github"
)

type labelCommand struct {
	cli.Command `name:"label" short-desc:"Manage GitHub Issues labels."`
	githubOptions
	Repositories []string `short:"r" long:"repositories" env:"JANITOR_REPOSITORIES" env-delim:"," required:"true" description:"Repositories to apply to, can be specified as user/name or as user/ for all repositories under the same account."`
	Name         string   `short:"n" long:"name" required:"true" desc:"Label name"`
	Description  string   `short:"d" long:"description" description:"Label description"`
	Color        string   `short:"c" long:"color" description:"Label color"`
	From         []string `long:"from" description:"Rename from old names, if they exist."`
}

func (c labelCommand) ExecuteContext(ctx context.Context, args []string) error {
	client := c.newGithubClient()
	repos, err := client.List(ctx, c.Repositories)
	if err != nil {
		return err
	}

	label := &github.Label{
		Name:        c.Name,
		Color:       c.Color,
		Description: c.Description,
		From:        c.From,
	}

	for _, repo := range repos {
		err := client.AddLabel(ctx, repo, label)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	app.AddCommand(&labelCommand{})
}
