package util

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(token, baseURL string) (*github.Client, error) {
	var tc *http.Client
	if token != "" {
		tc = oauth2.NewClient(
			context.Background(),
			oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: token,
			}),
		)
	}

	client := github.NewClient(tc)

	if baseURL != "" {
		u, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		client.BaseURL = u
	}

	return client, nil
}

// TODO: to Usecase, Repository
func GetAllRepository(client *github.Client, org string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	ctx := context.Background()
	option := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, option)
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
