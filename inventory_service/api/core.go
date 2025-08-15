package api

import (
	"context"

	"github.com/kripst/krosovka/inventory_service/internal/storage"
	pb "github.com/kripst/krosovka/inventory_service/proto"
	"go.uber.org/zap"
)

type ApiServer interface {
	CreateSneakers(ctx context.Context, in *pb.CreateSneakersRequest) (*pb.Response, error)
	GetSneakers(ctx context.Context, in *pb.GetSneakersRequest) (*pb.GetSneakersResponse, error)
	UpdateSneakers(ctx context.Context, in *pb.UpdateSneakersRequest) (*pb.Response, error)
	DeleteSneakers(ctx context.Context, in *pb.DeleteSneakersRequest) (*pb.Response, error)
}

type ApiServerImpl struct {
	pb.UnimplementedInventoryServiceServer
	s storage.Storage
	log *zap.Logger
}