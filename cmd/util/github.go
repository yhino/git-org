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
