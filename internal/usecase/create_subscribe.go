package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"subscribes/internal/domain/model"
	externalRef0 "subscribes/openapi"
	"subscribes/openapi/subscribe"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
)

const (
	emptyUuid = "00000000-0000-0000-0000-000000000000"
)

func (u UseCase) CreateSubscribe(ctx context.Context, request subscribe.CreateSubscribeRequestObject) (subscribe.CreateSubscribeResponseObject, error) {
	if vErr := validateCreateRequest(request.Body); vErr != nil {
		return vErr, nil
	}

	var endDate = mustParsePDate(request.Body.EndDate)

	// add month to start_date if end_date is nil
	if endDate == nil {
		endDate = new(mustParseDate(request.Body.StartDate).AddDate(0, 1, 0))
	}

	s, err := u.repo.Create(ctx,
		model.Subscribe{
			UserId:      model.ID(request.Body.UserId),
			ServiceName: model.ServiceName(request.Body.ServiceName),
			Price:       model.Price(request.Body.Price),
			StartDate:   model.Date(mustParseDate(request.Body.StartDate)),
			EndDate:     (*model.Date)(endDate),
		})

	if err != nil {
		slog.ErrorContext(ctx, "subscribe creating is failed",
			"status", 500, "message", err.Error())

		return subscribe.CreateSubscribe500ApplicationProblemPlusJSONResponse{
			Error: externalRef0.BaseError{
				Code:    externalRef0.Internal,
				Message: err.Error(),
				Params:  nil,
			},
		}, nil
	}

	slog.InfoContext(ctx, fmt.Sprintf("subscribe with id %s is created", uuid.UUID(s.ID).String()))

	return subscribe.CreateSubscribe201JSONResponse{
		GetSubscribeResponseJSONResponse: subscribe.GetSubscribeResponseJSONResponse{
			EndDate:     convertPDate(s.EndDate),
			Id:          openapitypes.UUID(s.ID),
			Price:       request.Body.Price,
			ServiceName: request.Body.ServiceName,
			StartDate:   request.Body.StartDate,
			UserId:      request.Body.UserId,
		},
	}, nil
}

func mustParseDate(date string) time.Time {
	pTime, _ := time.Parse("01-2006", date)
	return pTime
}

func mustParsePDate(date *string) *time.Time {
	if date == nil {
		return nil
	}
	pTime, _ := time.Parse("01-2006", *date)
	return new(pTime)
}

func validateCreateRequest(body *subscribe.CreateSubscribeJSONRequestBody) subscribe.CreateSubscribeResponseObject {
	if body.UserId.String() == emptyUuid {
		slog.Error("Need user id")

		return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "Need user id",
					Params:  nil,
				},
			},
		}
	}

	if body.Price == 0 {
		slog.Error("Need price")

		return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "Need price",
					Params:  nil,
				},
			},
		}
	}

	if body.Price < 0 {
		slog.Error("Price must be greater than zero")

		return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "Price must be greater than zero",
					Params:  nil,
				},
			},
		}
	}

	if body.ServiceName == "" {
		slog.Error("Need service name")

		return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "Need service name",
					Params:  nil,
				},
			},
		}
	}

	if utf8.RuneCountInString(body.ServiceName) > 255 {
		slog.Error("Service name must be less or equal than 255 characters")

		return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "Service name must be less or equal than 255 characters",
					Params:  nil,
				},
			},
		}
	}

	if body.StartDate == "" {
		slog.Error("Need start date")

		return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "Need start date",
					Params:  nil,
				},
			},
		}
	}

	stDate, err := time.Parse("01-2006", body.StartDate)

	if err != nil {
		slog.Error("Start date format error. Need MM-YYYY")

		return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
			ApiErrorResponse: externalRef0.ApiErrorResponse{
				Error: externalRef0.BaseError{
					Code:    externalRef0.BadRequest,
					Message: "Start date format error. Need MM-YYYY",
					Params:  nil,
				},
			},
		}
	}

	if body.EndDate != nil && *body.EndDate != "" {
		endDate, err := time.Parse("01-2006", *body.EndDate)

		if err != nil {
			slog.Error("End date format error. Need MM-YYYY")

			return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
				ApiErrorResponse: externalRef0.ApiErrorResponse{
					Error: externalRef0.BaseError{
						Code:    externalRef0.BadRequest,
						Message: "End date format error. Need MM-YYYY",
						Params:  nil,
					},
				},
			}
		}

		if !endDate.After(stDate) {
			slog.Error("End date error. Need be after start date")

			return subscribe.CreateSubscribe400ApplicationProblemPlusJSONResponse{
				ApiErrorResponse: externalRef0.ApiErrorResponse{
					Error: externalRef0.BaseError{
						Code:    externalRef0.BadRequest,
						Message: "End date error. Need be after start date",
						Params:  nil,
					},
				},
			}
		}
	}

	return nil
}
