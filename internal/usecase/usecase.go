package usecase

import (
	"context"
	"subscribes/internal/domain/model"
)

type ARepo interface {
	Get(id model.ID) (model.Subscribe, error)
	List(ctx context.Context, pageSize int, page int) ([]model.Subscribe, error)
	Delete(id model.ID) error
	Create(ctx context.Context, subscribe model.Subscribe) (model.Subscribe, error)
	Save(subscribe model.Subscribe) error
	PriceStat(userId *model.ID, name *model.ServiceName, startDate *model.Date, endDate *model.Date) (int, error)
}

type UseCase struct {
	repo ARepo
}

func NewUseCase(repo ARepo) *UseCase {
	return &UseCase{repo: repo}
}
