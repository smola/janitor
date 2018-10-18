package main

import (
	"context"
	"net/http"

	xgithub "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-cli.v0"

	"github.com/smola/janitor/github"
)

var (
	version string
	build   string
)

var app = cli.New("janitor", version, build, "Repository mainteinance tasks.")

type githubOptions struct {
	GithubToken string `long:"github-token" env:"GITHUB_TOKEN" description:"GitHub token."`
}

func (o githubOptions) newGithubClient() *github.Client {
	var httpClient *http.Client
	if o.GithubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: o.GithubToken},
		)
		httpClient = oauth2.NewClient(context.Background(), ts)
	}

	return &github.Client{
		Client: xgithub.NewClient(httpClient),
	}
}

func main() {
	app.RunMain()
}
