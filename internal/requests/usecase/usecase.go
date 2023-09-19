package usecase

import "github.com/SlavaShagalov/proxy-server/internal/requests"

type usecase struct {
	repo requests.Repository
}

func New(repo requests.Repository) requests.Usecase {
	return &usecase{repo: repo}
}

func (uc *usecase) Create(params *requests.Request) error {
	return uc.repo.Create(params)
}

func (uc *usecase) Get(id string) (*requests.Request, error) {
	return uc.repo.Get(id)
}

func (uc *usecase) List() ([]requests.Request, error) {
	return uc.repo.List()
}
