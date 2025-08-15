package storage

import (
	"context"

	"github.com/kripst/krosovka/inventory_service/internal/model"
)

type Storage interface {
	CreateSneakers(ctx context.Context, sneakers []*model.Sneaker) error
	UpdateSneakers(ctx context.Context, sneakers []*model.Sneaker) error
	DeleteSneakers(ctx context.Context, sneakerIDs []int32) error
	GetSneakers(ctx context.Context) ([]*model.Sneaker, error)
	Close() error
}
