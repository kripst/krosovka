package api

import (
	"context"
	"net/http"
	"time"

	pb "github.com/kripst/krosovka/inventory_service/proto"
	"go.uber.org/zap"
)


func(a *ApiServerImpl) DeleteSneakers(ctx context.Context, in *pb.DeleteSneakersRequest) (*pb.Response, error) {
	sneakerIDs := in.GetSneakerId()

	response := &pb.Response{}
	response.RequestId = in.GetRequestId()
	response.StatusCode = http.StatusOK
	response.SneakerIds = in.GetSneakerId()
	response.Timestamp = time.Now().String()

	if err := a.s.DeleteSneakers(ctx, sneakerIDs); err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.ErrorMessage = err.Error()
		response.Status = 1 // FAILURE
		response.SneakerIds = nil

		return response, err
	}

	a.log.Info("sneakers soft deleted", zap.Any("sneakerIDs", response.SneakerIds))
	return response, nil
}