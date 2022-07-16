package products

type Service struct {
	database Database
	cache    Cache
	notifier Notifier
}

func NewService(database Database, cache Cache, notifier Notifier) *Service {
	return &Service{database: database, cache: cache, notifier: notifier}
}

func (s Service) GetProductById(id string) (Product, error) {
	found, cachedProduct, cacheGetErr := s.cache.GetProductById(id)
	if cacheGetErr != nil {
		s.notifier.NotifyError(cacheFailedError)
		found = false
	}
	if found == true {
		return cachedProduct, nil
	}

	dbProduct, dbGetErr := s.database.GetProductById(id)
	if dbGetErr != nil {
		return Product{}, dbGetErr
	}

	cacheSetErr := s.cache.SetProduct(dbProduct)
	if cacheSetErr != nil {
		s.notifier.NotifyError(cacheFailedError)
	}

	return dbProduct, nil
}
