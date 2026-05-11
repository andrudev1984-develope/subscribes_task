package usecase

import (
	"context"
	"subscribes/internal/domain/model"
	externalRef0 "subscribes/openapi"
	"subscribes/openapi/subscribe"
	"time"
)

func (u UseCase) GetCalcedSubscribes(ctx context.Context, request subscribe.GetCalcedSubscribesRequestObject) (subscribe.GetCalcedSubscribesResponseObject, error) {
	if vErr := validateGCRequest(request); vErr != nil {
		return vErr, nil
	}

	pStat, err := u.repo.PriceStat(
		(*model.ID)(request.Params.UserId),
		(*model.ServiceName)(request.Params.ServiceName),
		(*model.Date)(mustParsePDate(request.Params.StartDate)),
		(*model.Date)(mustParsePDate(request.Params.EndDate)))

	if err != nil {
		return subscribe.GetCalcedSubscribes500ApplicationProblemPlusJSONResponse{
			Error: externalRef0.BaseError{
				Code:    externalRef0.Internal,
				Message: err.Error(),
				Params:  nil,
			},
		}, nil
	}

	return subscribe.GetCalcedSubscribes200JSONResponse{
		GetSubscribesPriceStatJSONResponse: subscribe.GetSubscribesPriceStatJSONResponse{
			TotalPrice: pStat,
		},
	}, nil
}

func validateGCRequest(request subscribe.GetCalcedSubscribesRequestObject) subscribe.GetCalcedSubscribesResponseObject {
	if request.Params.StartDate != nil {
		_, err := time.Parse("01-2006", *request.Params.StartDate)
		if err != nil {
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
