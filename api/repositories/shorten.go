package repositories

type ShortenRepository struct {
}

func NewShortenRepository() *ShortenRepository {
	return &ShortenRepository{}
}

func (repo ShortenRepository) GetAll() (err error) {

	return nil
}
