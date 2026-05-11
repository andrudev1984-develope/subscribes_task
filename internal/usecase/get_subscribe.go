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
	"time"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

func (u UseCase) GetSubscribe(ctx context.Context, request subscribe.GetSubscribeRequestObject) (subscribe.GetSubscribeResponseObject, error) {
	var s, err = u.repo.Get(model.ID(request.Id))

	if err != nil {
		var sErr, ok = errors.AsType[*out.SubscribeError](err)

		if ok {
			slog.ErrorContext(ctx, fmt.Sprintf("subscribe with id %s getting problem", request.Id.String()),
				"status", sErr.Code, "message", sErr.Error())

			return subscribe.GetSubscribe404ApplicationProblemPlusJSONResponse{
				Error: externalRef0.BaseError{
					Code:    "404",
					Message: sErr.Error(),
					Params:  nil,
				},
			}, nil
		}

		slog.ErrorContext(ctx, fmt.Sprintf("subscribe with id %s getting problem", request.Id.String()),
			"status", 500, "message", err.Error())

		return nil, err
	}

	slog.InfoContext(ctx, fmt.Sprintf("subscribe with id %s is received", request.Id.String()))

	return subscribe.GetSubscribe200JSONResponse{
		GetSubscribeResponseJSONResponse: subscribe.GetSubscribeResponseJSONResponse{
			EndDate:     convertPDate(s.EndDate),
			Id:          openapitypes.UUID(s.ID),
			Price:       int(s.Price),
			ServiceName: string(s.ServiceName),
			StartDate:   convertDate(s.StartDate),
			UserId:      openapitypes.UUID(s.UserId),
		},
	}, nil
}

func convertDate(date model.Date) string {
	return time.Time(date).Format("01-2006")
}

func convertPDate(date *model.Date) *string {
	if date == nil {
		return nil
	}

	return new(convertDate(*date))
}
