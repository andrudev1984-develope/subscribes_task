package usecase

import (
	"context"
	"log/slog"
	externalRef0 "subscribes/openapi"
	"subscribes/openapi/subscribe"

	openapitypes "github.com/oapi-codegen/runtime/types"
)

func (u UseCase) GetSubscribesList(ctx context.Context, request subscribe.GetSubscribesListRequestObject) (subscribe.GetSubscribesListResponseObject, error) {
	if eResp := validateListRequest(request); eResp != nil {
		return eResp, nil
	}

	var pageSize int
	var pageNumber int

	if request.Params.PageSize == nil || *request.Params.PageSize == 0 {
		pageSize = 50
	} else {
		pageSize = *request.Params.PageSize
	}

	if request.Params.PageNumber == nil || *request.Params.PageNumber == 0 {
		pageNumber = 1
	} else {
		pageNumber = *request.Params.PageNumber
	}

	subscribes, err := u.repo.List(ctx, pageSize, pageNumber)

	defer func() {
		subscribes = nil
	}()

	if err != nil {
		slog.ErrorContext(ctx, "subscribes list getting problem",
			"status", 500, "message", err.Error())

		return subscribe.GetSubscribesList500ApplicationProblemPlusJSONResponse{
			Error: externalRef0.BaseError{
				Code:    externalRef0.Internal,
				Message: err.Error(),
				Params:  nil,
			},
		}, nil
	}

	var subscribesOut = make([]subscribe.SubscribeOuInfo, len(subscribes))

	for i, s := range subscribes {
		subscribesOut[i] =
			subscribe.SubscribeOuInfo{
				EndDate:     convertPDate(s.EndDate),
				Id:          openapitypes.UUID(s.ID),
				Price:       int(s.Price),
				ServiceName: string(s.ServiceName),
				StartDate:   convertDate(s.StartDate),
				UserId:      openapitypes.UUID(s.UserId),
			}
	}

	slog.InfoContext(ctx, "subscribes list getting success")

	return subscribe.GetSubscribesList200JSONResponse(subscribesOut), nil
}

func validateListRequest(request subscribe.GetSubscribesListRequestObject) subscribe.GetSubscribesListResponseObject {
	if request.Params.PageSize != nil && *request.Params.PageSize <= 0 {
		slog.Error("Page size must be greater than zero")

		return subscribe.GetSubscribesList400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "page size must be greater than zero",
					Params:  nil,
				},
			},
		}
	}

	if request.Params.PageNumber != nil && *request.Params.PageNumber <= 0 {
		slog.Error("Page size must be greater than zero")

		return subscribe.GetSubscribesList400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "page number must be greater than zero",
					Params:  nil,
				},
			},
		}
	}

	return nil
}
