package services

import (
	"context"

	"github.com/ChrisMarSilva/cms-url-shortener/entities"
	opentracing "github.com/opentracing/opentracing-go"
)

type DefaultService struct {
}

func NewDefaultService() *DefaultService {
	return &DefaultService{}
}

func (service *DefaultService) NotFound(ctx context.Context) *entities.HttpResponse {

	sp, _ := opentracing.StartSpanFromContext(ctx, "DefaultService.NotFound")
	defer sp.Finish()

	return entities.NotFound("")
}
