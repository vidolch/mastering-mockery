package products

type Cache interface {
	GetProductById(id string) (bool, Product, error)
	SetProduct(product Product) error
}

var cacheProducts = map[string]Product{}

type CacheImpl struct{}

func NewCache() *CacheImpl {
	return &CacheImpl{}
}

func (c CacheImpl) GetProductById(id string) (bool, Product, error) {
	product, ok := cacheProducts[id]

	return ok, product, nil
}

func (c CacheImpl) SetProduct(product Product) error {
	cachedProduct := product
	cachedProduct.Name = cachedProduct.Name + " - From Cache"
	cacheProducts[product.ID] = cachedProduct

	return nil
}
