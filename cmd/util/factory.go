package util

import (
	"os"

	"github.com/google/go-github/github"
	"github.com/yhinoz/git-org/service"
)

type Factory interface {
	GithubClient() (*github.Client, error)
	RepositoryService() (service.RepositoryService, error)
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

func (f *factoryImpl) RepositoryService() (service.RepositoryService, error) {
	githubClient, err := f.GithubClient()
	if err != nil {
		return nil, err
	}

	service := service.NewRepositoryService(
		githubClient,
	)

	return service, nil
}
