package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Estriper0/subscription_service/internal/repository"
	"github.com/Estriper0/subscription_service/internal/repository/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepo struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepo(db *pgxpool.Pool) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: db,
	}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *models.SubscriptionCreate) (int, error) {
	query := `
		INSERT INTO subscription (service_name, price, user_id, start_date, end_date) 
			VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	var id int

	err := r.db.QueryRow(ctx, query, s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate).Scan(&id)
	return id, err
}

func (r *SubscriptionRepo) GetById(ctx context.Context, id int) (*models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscription 
		WHERE id = $1
	`
	var subscription models.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&subscription.Id,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserId,
		&subscription.StartDate,
		&subscription.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("db:SubscriptionRepo.GetById:QueryRow - %s", err.Error())
	}

	return &subscription, nil
}

func (r *SubscriptionRepo) GetByUser(ctx context.Context, userId uuid.UUID) ([]*models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscription 
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("db:SubscriptionRepo.GetByUser:Query - %s", err.Error())
	}

	var subscriptions []*models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&subscription.Id,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserId,
			&subscription.StartDate,
			&subscription.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf("db:SubscriptionRepo.GetByUser:Scan - %s", err.Error())
		}
		subscriptions = append(subscriptions, &subscription)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepo) DeleteById(ctx context.Context, id int) (*models.Subscription, error) {
	query := `
		DELETE FROM subscription 
			WHERE id = $1
		RETURNING id, service_name, price, user_id, start_date, end_date
	`
	var subscription models.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&subscription.Id,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserId,
		&subscription.StartDate,
		&subscription.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("db:SubscriptionRepo.DeleteById:QueryRow - %s", err.Error())
	}

	return &subscription, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, s *models.SubscriptionUpdate) (*models.Subscription, error) {
	query := `
		UPDATE subscription 
		SET
			service_name = COALESCE($1, service_name),
			price = COALESCE($2, price),
			start_date = COALESCE($3, start_date),
			end_date = COALESCE($4, end_date)
		WHERE
			id = $5
		RETURNING id, service_name, price, user_id, start_date, end_date
	`
	var subscription models.Subscription

	err := r.db.QueryRow(ctx, query, s.ServiceName, s.Price, s.StartDate, s.EndDate, s.Id).Scan(
		&subscription.Id,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserId,
		&subscription.StartDate,
		&subscription.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == repository.PgCodeConstrainError {
				return nil, repository.ErrIncorrectTime
			}
		}
		return nil, fmt.Errorf("db:SubscriptionRepo.Update:QueryRow - %s", err.Error())
	}

	return &subscription, nil
}

func (r *SubscriptionRepo) GetPriceByFilter(ctx context.Context, user_id *uuid.UUID, service_name *string, start_date, end_date time.Time) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0) 
        	FROM subscription
        WHERE start_date <= $1 
			AND 
			(
				end_date IS NULL OR
				end_date >= $2
			) 
    `
	c := 3
	args := []interface{}{end_date, start_date}
	if user_id != nil {
		query += fmt.Sprintf(" AND user_id = $%d", c)
		c++
		args = append(args, *user_id)
	}
	if service_name != nil {
		query += fmt.Sprintf(" AND service_name = $%d", c)
		args = append(args, *service_name)
	}

	var total int
	err := r.db.QueryRow(ctx, query, args...).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("db:SubscriptionRepo.GetPriceByTime:QueryRow - %s", err.Error())
	}

	return total, nil
}
