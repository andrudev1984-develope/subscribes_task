package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"subscribes/internal/domain/model"
	"subscribes/internal/dto/out"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	getSingle = `SELECT id, user_id, service_name, price, start_date, end_date
					FROM public.subscribe WHERE id = $1`
	insertSingle = `INSERT INTO public.subscribe (user_id, service_name, price, start_date, end_date) 
					VALUES ($1, $2, $3, $4, $5) RETURNING id`
	deleteSingle = `DELETE FROM public.subscribe WHERE id = $1`
	updateSingle = `UPDATE public.subscribe SET user_id = $1, service_name = $2, 
                            price = $3, start_date = $4, end_date = $5 WHERE id = $6`
	getList = `SELECT id, user_id, service_name, price, start_date, end_date 
					FROM public.subscribe ORDER BY user_id, service_name ASC OFFSET $1 LIMIT $2`
	getPriceStat = `SELECT sum(s.price) FROM public.subscribe s 
                    	WHERE (($2::text = '') OR (s.service_name = $2)) and 
							  ($3::text = '' or s.start_date >= $3::date) and
							  ($4::text = '' or s.start_date <= $4::date) and
					          ($3::text = '' or s.end_date >= $3::date) and
							  ($4::text = '' or s.end_date <= $4::date) and
							  ($1::text = '' or s.user_id = $1::uuid)`
)

type Pgrep struct {
	db *pgxpool.Pool
}

func (p Pgrep) PriceStat(userId *model.ID, name *model.ServiceName, startDate *model.Date, endDate *model.Date) (int, error) {
	var pStat = sql.NullInt64{
		Int64: 0,
		Valid: false,
	}

	var cUserId string
	var cServiceName string
	var cStartDate string
	var cEndDate string

	if userId != nil {
		cUserId = uuid.UUID(*userId).String()
	}

	if name != nil {
		cServiceName = string(*name)
	}

	if startDate != nil {
		cStartDate = time.Time(*startDate).Format(time.DateOnly)
	}

	if endDate != nil {
		cEndDate = time.Time(*endDate).Format(time.DateOnly)
	}

	err := p.db.QueryRow(context.Background(), getPriceStat,
		cUserId,
		cServiceName,
		cStartDate,
		cEndDate).Scan(&pStat)

	if err != nil {
		return 0, nil
	}

	return int(pStat.Int64), err
}

func (p Pgrep) List(ctx context.Context, pageSize int, page int) ([]model.Subscribe, error) {
	res, err := p.db.Query(ctx, getList, (page-1)*pageSize, pageSize)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	defer res.Close()

	var subscribes = make([]model.Subscribe, 0, pageSize)

	for res.Next() {
		var s model.Subscribe
		var stDate pgtype.Date
		var enDate pgtype.Date

		err = res.Scan(&s.ID, &s.UserId, &s.ServiceName, &s.Price, &stDate, &enDate)

		if err != nil {
			return nil, fmt.Errorf("database read error: %w", err)
		}

		s.StartDate = model.Date(stDate.Time)

		if enDate.Valid {
			s.EndDate = (*model.Date)(&enDate.Time)
		}

		subscribes = append(subscribes, s)
	}

	return subscribes, nil
}

func NewPgrep(db *pgxpool.Pool) *Pgrep {
	return &Pgrep{db: db}
}

func (p Pgrep) Create(ctx context.Context, subscribe model.Subscribe) (model.Subscribe, error) {
	var crId uuid.UUID

	err := p.db.QueryRow(ctx, insertSingle,
		subscribe.UserId,
		subscribe.ServiceName,
		subscribe.Price,
		pgtype.Date{
			Time:             time.Time(subscribe.StartDate),
			InfinityModifier: 0,
			Valid:            true,
		},
		pgtype.Date{
			Time:             time.Time(*subscribe.EndDate),
			InfinityModifier: 0,
			Valid:            true,
		}).
		Scan(&crId)

	if err != nil {
		return model.Subscribe{}, fmt.Errorf("database error: %w", err)
	}

	return model.Subscribe{ID: model.ID(crId)}, nil
}

func (p Pgrep) Get(id model.ID) (model.Subscribe, error) {
	var subs model.Subscribe

	var stDate pgtype.Date
	var enDate pgtype.Date

	err := p.db.QueryRow(context.Background(), getSingle, id).
		Scan(&subs.ID, &subs.UserId, &subs.ServiceName, &subs.Price, &stDate, &enDate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscribe{}, &out.SubscribeError{
				Code:    404,
				Details: fmt.Sprintf("subscribe with id %s is not found", uuid.UUID(id).String()),
			}
		}
		return model.Subscribe{}, err
	}

	subs.StartDate = model.Date(stDate.Time)

	if enDate.Valid {
		subs.EndDate = (*model.Date)(&enDate.Time)
	}

	return subs, nil
}

func (p Pgrep) Save(subscribe model.Subscribe) error {
	res, err := p.db.Exec(context.Background(), updateSingle,
		subscribe.UserId, subscribe.ServiceName, subscribe.Price,
		pgtype.Date{
			Time:             time.Time(subscribe.StartDate),
			InfinityModifier: 0,
			Valid:            true,
		}, pgtype.Date{
			Time:             time.Time(*subscribe.EndDate),
			InfinityModifier: 0,
			Valid:            true,
		}, subscribe.ID)

	if err != nil {
		return fmt.Errorf("database error %d: %w", subscribe.ID, err)
	}

	if res.RowsAffected() == 0 {
		return &out.SubscribeError{
			Code:    404,
			Details: fmt.Sprintf("subscribe with id %s is not found", uuid.UUID(subscribe.ID).String())}
	}

	return nil
}

func (p Pgrep) Delete(id model.ID) error {
	res, err := p.db.Exec(context.Background(), deleteSingle, id)

	if err != nil {
		return fmt.Errorf("database error %d: %w", id, err)
	}

	if res.RowsAffected() == 0 {
		return &out.SubscribeError{
			Code:    404,
			Details: fmt.Sprintf("subscribe with id %s is not found", uuid.UUID(id).String())}
	}

	return nil
}
