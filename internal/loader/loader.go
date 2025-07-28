package loader

type ServiceLoader interface {
	LoadProducts(path string) error
	LoadPrices(path string) error
}
