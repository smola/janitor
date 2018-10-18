package github

import (
	"context"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-log.v1"
)

type Label struct {
	Name        string
	Description string
	Color       string
	From        []string
}

func (l Label) toRemote() *github.Label {
	rlabel := &github.Label{
		Name: &l.Name,
	}

	if len(l.Color) != 0 {
		rlabel.Color = &l.Color
	}

	if len(l.Description) != 0 {
		rlabel.Description = &l.Description
	}

	return rlabel
}

func (c *Client) AddLabel(ctx context.Context, repo *Repository, label *Label) error {
	opts := &github.ListOptions{}
	labels, resp, err := c.Client.Issues.ListLabels(ctx, repo.Owner, repo.Name, opts)
	if err != nil {
		log.Debugf("unexpected error, response was: %#v", resp)
		return err
	}

	for _, rlabel := range labels {
		if *rlabel.Name == label.Name {
			return c.editLabel(ctx, repo, label.Name, label)
		}

		for _, from := range label.From {
			if *rlabel.Name == from {
				return c.editLabel(ctx, repo, from, label)
			}
		}
	}

	return c.createLabel(ctx, repo, label)
}

func (c *Client) createLabel(ctx context.Context, repo *Repository, label *Label) error {
	_, resp, err := c.Client.Issues.CreateLabel(ctx, repo.Owner, repo.Name, label.toRemote())
	if err != nil {
		log.Debugf("unexpected error, response was: %#v", resp)
		return err
	}

	return nil
}

func (c *Client) editLabel(ctx context.Context, repo *Repository, old string, label *Label) error {
	_, resp, err := c.Client.Issues.EditLabel(ctx, repo.Owner, repo.Name, old, label.toRemote())
	if err != nil {
		log.Debugf("unexpected error, response was: %#v", resp)
		return err
	}

	return nil
}
