package repositories

import (
	"time"
)

type URLRepository struct {
}

func NewURLRepository() *URLRepository {
	return &URLRepository{}
}

func (repo URLRepository) GetById(_id string) (value string, err error) {
	// value, err := r.Get(databases.Ctx, _id).Result()
	return
}

func (repo URLRepository) Create(_id string, url string, expiry time.Duration) (err error) {
	// err = r.Set(databases.Ctx, _id, url, expiry).Err()
	return
}
