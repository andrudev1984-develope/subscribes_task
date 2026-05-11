package usecase

import (
	"context"
	"log/slog"
	"subscribes/internal/domain/model"
	externalRef0 "subscribes/openapi"
	"subscribes/openapi/subscribe"
	"time"
)

func (u UseCase) GetCalcedSubscribes(ctx context.Context, request subscribe.GetCalcedSubscribesRequestObject) (subscribe.GetCalcedSubscribesResponseObject, error) {
	if vErr := validateStatPricesRequest(request); vErr != nil {
		return vErr, nil
	}

	pStat, err := u.repo.PriceStat(ctx,
		(*model.ID)(request.Params.UserId),
		(*model.ServiceName)(request.Params.ServiceName),
		(*model.Date)(mustParsePDate(request.Params.StartDate)),
		(*model.Date)(mustParsePDate(request.Params.EndDate)))

	if err != nil {
		slog.ErrorContext(ctx, "getting subscribe prices is failed",
			"status", 500, "message", err.Error())

		return subscribe.GetCalcedSubscribes500ApplicationProblemPlusJSONResponse{
			Error: externalRef0.BaseError{
				Code:    externalRef0.Internal,
				Message: err.Error(),
				Params:  nil,
			},
		}, nil
	}

	slog.InfoContext(ctx, "getting subscribe prices is successes")

	return subscribe.GetCalcedSubscribes200JSONResponse{
		GetSubscribesPriceStatJSONResponse: subscribe.GetSubscribesPriceStatJSONResponse{
			TotalPrice: pStat,
		},
	}, nil
}

func validateStatPricesRequest(request subscribe.GetCalcedSubscribesRequestObject) subscribe.GetCalcedSubscribesResponseObject {
	if request.Params.StartDate != nil {
		_, err := time.Parse("01-2006", *request.Params.StartDate)

		if err != nil {
			slog.Error("Start date format error. Need MM-YYYY")

			return subscribe.GetCalcedSubscribes400ApplicationProblemPlusJSONResponse{
				ApiErrorResponse: externalRef0.ApiErrorResponse{
					Error: externalRef0.BaseError{
						Code:    externalRef0.BadRequest,
						Message: "Start date format error. Need MM-YYYY",
						Params:  nil,
					},
				},
			}
		}
	}

	if request.Params.EndDate != nil {
		_, err := time.Parse("01-2006", *request.Params.EndDate)

		if err != nil {
			slog.Error("End date format error. Need MM-YYYY")

			return subscribe.GetCalcedSubscribes400ApplicationProblemPlusJSONResponse{
				ApiErrorResponse: externalRef0.ApiErrorResponse{
					Error: externalRef0.BaseError{
						Code:    externalRef0.BadRequest,
						Message: "End date format error. Need MM-YYYY",
						Params:  nil,
					},
				},
			}
		}
	}

	return nil
}
