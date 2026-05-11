package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"subscribes/internal/domain/model"
	"subscribes/internal/dto/out"
	externalRef0 "subscribes/openapi"
	"subscribes/openapi/subscribe"
)

func (u UseCase) DeleteSubscribe(ctx context.Context, request subscribe.DeleteSubscribeRequestObject) (subscribe.DeleteSubscribeResponseObject, error) {
	var err = u.repo.Delete(ctx, model.ID(request.Id))

	if err != nil {
		var sErr, ok = errors.AsType[*out.SubscribeError](err)

		if ok {
			slog.ErrorContext(ctx, fmt.Sprintf("subscribe with id %s delete problem", request.Id.String()),
				"status", sErr.Code, "message", sErr.Error())

			return subscribe.DeleteSubscribe404ApplicationProblemPlusJSONResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.NotFound,
					Message: sErr.Error(),
					Params:  nil,
				},
			}, nil
		}

		slog.ErrorContext(ctx, fmt.Sprintf("subscribe with id %s delete problem", request.Id.String()),
			"status", 500, "message", err.Error())

		return subscribe.DeleteSubscribe500ApplicationProblemPlusJSONResponse{
			Error: externalRef0.BaseError{
				Code:    externalRef0.Internal,
				Message: err.Error(),
				Params:  nil,
			},
		}, nil
	}

	slog.InfoContext(ctx, fmt.Sprintf("subscribe with id %s is deleted", request.Id.String()))

	return subscribe.DeleteSubscribe204Response{}, nil
}
