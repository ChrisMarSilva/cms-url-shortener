package services

import (
	"github.com/ChrisMarSilva/cms-url-shortener/repositories"
)

type ResolveService struct {
	repo repositories.ResolveRepository
}

func NewResolveService(repo repositories.ResolveRepository) *ResolveService {
	return &ResolveService{
		repo: repo,
	}
}

func (service *ResolveService) ResolveURL(url string) (err error) {

	return nil
}
