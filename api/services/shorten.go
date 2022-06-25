package services

import (
	"github.com/ChrisMarSilva/cms-url-shortener/entities"
	"github.com/ChrisMarSilva/cms-url-shortener/repositories"
)

type ShortenService struct {
	repo repositories.ShortenRepository
}

func NewShortenService(repo repositories.ShortenRepository) *ShortenService {
	return &ShortenService{
		repo: repo,
	}
}

func (service *ShortenService) ShortenURL(body entities.Request) (err error) {

	return nil
}
