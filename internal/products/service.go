package products

import "context"

type Service interface {
	ListProducts(ctx context.Context) ([]string, error)
}

type svc struct {
}

func NewService() Service {
	return &svc{}
}

func (s *svc) ListProducts(ctx context.Context) ([]string, error) {
	return []string{"hello", "world"}, nil
}
