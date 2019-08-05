package github

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-log.v1"
)

type RepositoriesSpec struct {
	Owner string
	Names []string
}

type Repository struct {
	Owner     string
	Name      string
	Archived  bool
	Private   bool
	CreatedAt time.Time
	PushedAt  time.Time
	UpdatedAt time.Time
	// License as SPDXID.
	License     string
	Maintainers []*User
}

func (c *Client) List(ctx context.Context, repoMasks []string) ([]*Repository, error) {
	rms, err := parseRepoMasks(repoMasks)
	if err != nil {
		return nil, err
	}

	if len(rms) == 0 {
		return nil, fmt.Errorf("at least one repository mask is required")
	}

	var result []*Repository
	for _, rm := range rms {
		if len(rm.Names) == 0 {
			repos, err := c.listOrganization(ctx, rm.Owner)
			if err != nil {
				return result, err
			}

			result = append(result, repos...)
		} else {
			for _, name := range rm.Names {
				repo, err := c.getRepo(ctx, rm.Owner, name)
				if err != nil {
					return result, err
				}

				if repo == nil {
					log.Warningf("given repository does not exist: %s/%s", rm.Owner, name)
				} else {
					result = append(result, repo)
				}
			}
		}
	}

	return result, nil
}

func (c *Client) listOrganization(ctx context.Context, org string) ([]*Repository, error) {
	var result []*Repository
	opts := &github.RepositoryListOptions{}

	for {
		repos, resp, err := c.Client.Repositories.List(ctx, org, opts)
		if err != nil {
			log.Debugf("got unexpected error, respose was: %#v", resp)
			return nil, err
		}

		for _, repo := range repos {
			lrepo := repositoryRemoteToLocal(repo)
			lrepo.Owner = org
			result = append(result, lrepo)
		}

		if resp.NextPage > 0 {
			opts.Page = resp.NextPage
		} else {
			break
		}
	}

	return result, nil
}

func (c *Client) getRepo(ctx context.Context, org, repo string) (*Repository, error) {
	rrepo, resp, err := c.Client.Repositories.Get(ctx, org, repo)
	if err == nil {
		return repositoryRemoteToLocal(rrepo), nil
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	log.Debugf("got unexpected error, response was: %#v", resp)
	return nil, err
}

func parseRepoMasks(repoMasks []string) ([]*RepositoriesSpec, error) {
	sort.Strings(repoMasks)
	var result []*RepositoriesSpec
	for _, repoMask := range repoMasks {
		org, name, err := parseRepoMask(repoMask)
		if err != nil {
			return nil, err
		}

		var spec *RepositoriesSpec
		if len(result) == 0 || result[len(result)-1].Owner != org {
			spec = &RepositoriesSpec{
				Owner: org,
			}
			result = append(result, spec)
		} else {
			spec = result[len(result)-1]
		}

		if len(name) != 0 {
			spec.Names = append(spec.Names, name)
		}
	}

	return result, nil
}

func parseRepoMask(repoMask string) (string, string, error) {
	fields := strings.Split(repoMask, "/")
	if len(fields) != 2 {
		return "", "", fmt.Errorf("invalid repository mask: %s (use org/ or org/repo)", repoMask)
	}

	if len(fields[0]) == 0 {
		return "", "", fmt.Errorf("invalid repository mask: %s (organization cannot be empty)", repoMask)
	}

	return fields[0], fields[1], nil
}

func repositoryRemoteToLocal(rrepo *github.Repository) *Repository {
	if rrepo == nil {
		return nil
	}

	lrepo := &Repository{
		Name:      rrepo.GetName(),
		Archived:  rrepo.GetArchived(),
		Private:   rrepo.GetPrivate(),
		License:   rrepo.GetLicense().GetSPDXID(),
		CreatedAt: rrepo.GetCreatedAt().Time,
		PushedAt:  rrepo.GetPushedAt().Time,
		UpdatedAt: rrepo.GetUpdatedAt().Time,
	}

	if rrepo.Owner != nil {
		lrepo.Owner = rrepo.GetOwner().GetLogin()
	} else {
		lrepo.Owner = rrepo.GetOrganization().GetName()
	}

	return lrepo
}
