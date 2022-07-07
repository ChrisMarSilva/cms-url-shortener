package services

import (
	"context"

	"github.com/ChrisMarSilva/cms-url-shortener/databases"
	"github.com/ChrisMarSilva/cms-url-shortener/entities"
	"github.com/ChrisMarSilva/cms-url-shortener/repositories"
	"github.com/go-redis/redis/v8"
	opentracing "github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
)

type ResolveService struct {
	repo repositories.ResolveRepository
}

func NewResolveService(repo repositories.ResolveRepository) *ResolveService {
	return &ResolveService{
		repo: repo,
	}
}

func (service *ResolveService) ResolveURL(ctx context.Context, url string) *entities.HttpResponse {

	sp, _ := opentracing.StartSpanFromContext(ctx, "ResolveService.ResolveURL")
	defer sp.Finish()

	sp2 := opentracing.StartSpan("BD.Conn", opentracing.ChildOf(sp.Context()))
	r := databases.CreateClient(0)
	sp2.Finish()
	defer r.Close()

	sp3 := opentracing.StartSpan("BD.Get", opentracing.ChildOf(sp.Context()))
	value, err := r.Get(databases.Ctx, url).Result()
	sp3.Finish()
	if err == redis.Nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return entities.BadRequest("short not found in the database")
	} else if err != nil {
		sp.SetTag("error", true)
		sp.LogFields(spanlog.Error(err))
		return entities.BadRequest("cannot connext to DB")
		// return entities.ServerError() // "cannot connext to DB"
	}
	sp.SetTag("value", value)

	// rInr := databases.CreateClient(1)
	// defer rInr.Close()

	// _ = rInr.Incr(databases.Ctx, "counter")

	return entities.OkWithUrl("", value)
}
