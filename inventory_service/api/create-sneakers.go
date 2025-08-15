package api

import (
	"context"
	"net/http"
	"time"

	"github.com/kripst/krosovka/inventory_service/internal/model"
	pb "github.com/kripst/krosovka/inventory_service/proto"
	"go.uber.org/zap"
)

func (a *ApiServerImpl) CreateSneakers(ctx context.Context, in *pb.CreateSneakersRequest) (*pb.Response, error) {
	response := &pb.Response{}
	response.RequestId = in.GetRequestId()
	response.StatusCode = http.StatusOK
	response.SneakerIds = make([]int32, 0, len(in.GetSneakers()))
	response.Timestamp = time.Now().String()
	

	sneakers := make([]*model.Sneaker, 0, len(in.GetSneakers()))
	defer func() {
		if err := a.s.CreateSneakers(ctx, sneakers); err != nil {
			a.log.Error("ERROR: create sneakers", zap.Error(err))
		} else {
			a.log.Info("successfully created sneakers", zap.Int("quantity sneakers created", len(sneakers)))
		}
	}()

	for _, sneaker := range in.GetSneakers() {
		s := &model.Sneaker{}
		if err := s.FromGrpc(sneaker); err != nil {
			response.ErrorMessage = err.Error()
			response.StatusCode = http.StatusBadRequest
			response.Status = 2 // VALIDATION_ERROR
			a.log.Error("ERROR: bad request create sneakers", zap.Error(err))
			return response, err
		}
		sneakers = append(sneakers, s)
		response.SneakerIds = append(response.SneakerIds, s.ID)

	}

	return response, nil
}