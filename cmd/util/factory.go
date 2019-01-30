package util

import (
	"os"

	"github.com/google/go-github/github"
)

type Factory interface {
	GithubClient() (*github.Client, error)
}

type factoryImpl struct {
}

func NewFactory() Factory {
	return &factoryImpl{}
}

func (f *factoryImpl) GithubClient() (*github.Client, error) {
	return GetGithubClient(
		os.Getenv("GITHUB_ACCESS_TOKEN"),
		os.Getenv("GITHUB_API_BASE_URL"),
	)
}
