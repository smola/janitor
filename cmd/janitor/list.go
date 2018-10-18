package main

import (
	"context"
	"fmt"

	"github.com/smola/janitor/github"
	"gopkg.in/src-d/go-cli.v0"
)

type listCommand struct {
	cli.Command  `name:"list" short-desc:"List GitHub repositories."`
	Repositories []string `short:"r" long:"repositories" env:"JANITOR_REPOSITORIES" env-delim:"," required:"true" description:"Repositories to apply to, can be specified as user/name or as user/ for all repositories under the same account."`
}

func (c listCommand) ExecuteContext(ctx context.Context, args []string) error {
	client := github.Default
	repos, err := client.List(ctx, c.Repositories)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Printf("%s/%s\n", repo.Owner, repo.Name)
	}

	return nil
}

func init() {
	app.AddCommand(&listCommand{})
}
