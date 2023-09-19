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

//func (uc *usecase) FullUpdate(params *boards.FullUpdateParams) (*models.Board, error) {
//	return uc.repo.FullUpdate(params)
//}
//
//func (uc *usecase) PartialUpdate(params *boards.PartialUpdateParams) (*models.Board, error) {
//	return uc.repo.PartialUpdate(params)
//}
//
//func (uc *usecase) Delete(id string) error {
//	return uc.repo.Delete(id)
//}
