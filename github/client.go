package github

import (
	"github.com/google/go-github/github"
)

type Client struct {
	Client *github.Client
}

var Default = &Client{
	Client: github.NewClient(nil),
}
