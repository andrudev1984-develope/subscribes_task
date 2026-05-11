package usecase

import (
	"context"
	"errors"
	"log/slog"
	"subscribes/internal/domain/model"
	"subscribes/internal/dto/out"
	externalRef0 "subscribes/openapi"
	"subscribes/openapi/subscribe"
	"time"
)

func (u UseCase) UpdateSubscribe(ctx context.Context, request subscribe.UpdateSubscribeRequestObject) (subscribe.UpdateSubscribeResponseObject, error) {
	if vErr := validateUBody(request.Body); vErr != nil {
		return vErr, nil
	}

	var endDate = mustParsePDate(request.Body.EndDate)

	// add month to start_date if end_date is nil
	if endDate == nil {
		endDate = new(mustParseDate(request.Body.StartDate).AddDate(0, 1, 0))
	}

	err := u.repo.Save(ctx,
		model.Subscribe{
			ID:          model.ID(request.Id),
			UserId:      model.ID(request.Body.UserId),
			ServiceName: model.ServiceName(request.Body.ServiceName),
			Price:       model.Price(request.Body.Price),
			StartDate:   model.Date(mustParseDate(request.Body.StartDate)),
			EndDate:     (*model.Date)(endDate),
		})

	if err != nil {
		var sErr, ok = errors.AsType[*out.SubscribeError](err)

		if ok {
			return subscribe.UpdateSubscribe400ApplicationProblemPlusJSONResponse{
				ApiErrorResponse: externalRef0.ApiErrorResponse{
					Error: externalRef0.BaseError{
						Code:    "404",
						Message: sErr.Error(),
						Params:  nil,
					},
				},
			}, nil
		}

		return nil, err
	}

	return subscribe.UpdateSubscribe200JSONResponse{
		GetSubscribeResponseJSONResponse: subscribe.GetSubscribeResponseJSONResponse{
			EndDate:     new(convertDate(model.Date(*endDate))),
			Id:          request.Id,
			Price:       request.Body.Price,
			ServiceName: request.Body.ServiceName,
			StartDate:   request.Body.StartDate,
			UserId:      request.Body.UserId,
		},
	}, nil
}

func validateUBody(body *subscribe.CreateSubscribeJSONRequestBody) subscribe.UpdateSubscribeResponseObject {
	if body.UserId.String() == emptyUuid {
		slog.Error("Need user id")

		return subscribe.UpdateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    "400",
					Message: "Need user id",
					Params:  nil,
				},
			},
		}
	}

	if body.Price == 0 {
		slog.Error("Need price")

		return subscribe.UpdateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    "400",
					Message: "Need price",
					Params:  nil,
				},
			},
		}
	}

	if body.ServiceName == "" {
		slog.Error("Need service name")

		return subscribe.UpdateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    "400",
					Message: "Need service name",
					Params:  nil,
				},
			},
		}
	}

	if body.StartDate == "" {
		slog.Error("Need start date")

		return subscribe.UpdateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    "400",
					Message: "Need start date",
					Params:  nil,
				},
			},
		}
	}

	if _, err := time.Parse("01-2006", body.StartDate); err != nil {
		slog.Error("Start date format error. Need MM-YYYY")

		return subscribe.UpdateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    "400",
					Message: "Start date format error. Need MM-YYYY",
					Params:  nil,
				},
			},
		}
	}

	if body.EndDate != nil && *body.EndDate != "" {
		if _, err := time.Parse("01-2006", *body.EndDate); err != nil {
			slog.Error("End date format error. Need MM-YYYY")

			return subscribe.UpdateSubscribe400ApplicationProblemPlusJSONResponse{
				ApiErrorResponse: externalRef0.ApiErrorResponse{
					Error: externalRef0.BaseError{
						Code:    "400",
						Message: "End date format error. Need MM-YYYY",
						Params:  nil,
					},
				},
			}
		}
	}

	return nil
}
