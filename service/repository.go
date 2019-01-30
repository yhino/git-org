package service

import (
	"context"

	"github.com/google/go-github/github"
)

type RepositoryService interface {
	FetchAllByOrg(org string) ([]*github.Repository, error)
}

type repositoryService struct {
	github *github.Client
}

func NewRepositoryService(github *github.Client) RepositoryService {
	return &repositoryService{
		github: github,
	}
}

func (u *repositoryService) FetchAllByOrg(org string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	ctx := context.Background()
	option := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := u.github.Repositories.ListByOrg(ctx, org, option)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}

		option.ListOptions.Page = resp.NextPage
	}

	return allRepos, nil
}
