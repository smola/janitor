package main

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/src-d/go-cli.v0"
	log "gopkg.in/src-d/go-log.v1"

	"github.com/smola/janitor/github"
)

type listCommand struct {
	cli.Command `name:"list" short-desc:"List GitHub repositories."`
	githubOptions
	Repositories []string `short:"r" long:"repositories" env:"JANITOR_REPOSITORIES" env-delim:"," required:"true" description:"Repositories to apply to, can be specified as user/name or as user/ for all repositories under the same account."`
	Format       string   `long:"format" default:"{{ .Owner }}/{{ .Name }}" description:"template to use for results"`
}

func (c listCommand) ExecuteContext(ctx context.Context, args []string) error {
	tpl, err := template.New("format").Parse(c.Format)
	if err != nil {
		return fmt.Errorf("bad format: %s", err)
	}

	client := c.newGithubClient()
	repos, err := client.List(ctx, c.Repositories)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		repo.Maintainers, err = github.GetMaintainers(repo.Owner, repo.Name)
		if err != nil {
			log.Errorf(err, "retrieving maintainers for %s/%s", repo.Owner, repo.Name)
		}
		err := tpl.Execute(os.Stdout, repo)
		if err != nil {
			return err
		}

		fmt.Print("\n")
	}

	return nil
}

func init() {
	app.AddCommand(&listCommand{})
}
