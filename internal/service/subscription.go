package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Estriper0/subscription_service/internal/repository"
	"github.com/Estriper0/subscription_service/internal/repository/models"
	"github.com/Estriper0/subscription_service/internal/service/domain"
	"github.com/google/uuid"
)

type ISubscriptionRepo interface {
	Create(ctx context.Context, s *models.SubscriptionCreate) (int, error)
	GetById(ctx context.Context, id int) (*models.Subscription, error)
	GetByUser(ctx context.Context, userId uuid.UUID) ([]*models.Subscription, error)
	DeleteById(ctx context.Context, id int) (*models.Subscription, error)
	Update(ctx context.Context, s *models.SubscriptionUpdate) (*models.Subscription, error)
	GetPriceByFilter(ctx context.Context, user_id *uuid.UUID, service_name *string, start_date, end_date time.Time) (int, error)
}

type SubscriptionService struct {
	subscriptionRepo ISubscriptionRepo
	logger           *slog.Logger
}

func NewSubscriptionService(subscriptionRepo ISubscriptionRepo, logger *slog.Logger) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		logger:           logger,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, subscription *domain.SubscriptionCreate) (int, error) {
	startDate, _ := time.Parse("01-2006", subscription.StartDate)
	model := &models.SubscriptionCreate{
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserId:      subscription.UserId,
		StartDate:   startDate,
	}
	if subscription.EndDate != nil {
		endDate, _ := time.Parse("01-2006", *subscription.EndDate)
		if endDate.Before(startDate) {
			return 0, ErrIncorrectTime
		}
		model.EndDate = sql.NullTime{Time: endDate, Valid: true}
	}

	id, err := s.subscriptionRepo.Create(ctx, model)
	if err != nil {
		s.logger.Error("SubscriptionService.Add:subscriptionRepo.Create - Internal error", slog.String("error", err.Error()))
		return 0, ErrInternal
	}

	s.logger.Info(fmt.Sprintf("The subscription id=%d has been created", id))
	return id, err
}

func (s *SubscriptionService) GetByUser(ctx context.Context, user_id uuid.UUID) ([]*domain.Subscription, error) {
	models, err := s.subscriptionRepo.GetByUser(ctx, user_id)
	if err != nil {
		s.logger.Error("SubscriptionService.GetByUser:subscriptionRepo.GetByUser - Internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	var subscriptions []*domain.Subscription
	for _, m := range models {
		subscription := &domain.Subscription{
			Id:          m.Id,
			ServiceName: m.ServiceName,
			Price:       m.Price,
			UserId:      m.UserId,
			StartDate:   m.StartDate.Format("01-2006"),
		}

		if m.EndDate.Valid {
			subscription.EndDate = new(string)
			*subscription.EndDate = m.EndDate.Time.Format("01-2006")
		}
		subscriptions = append(subscriptions, subscription)
	}
	s.logger.Info("All user subscriptions were received successfully")

	return subscriptions, err
}

func (s *SubscriptionService) GetById(ctx context.Context, id int) (*domain.Subscription, error) {
	model, err := s.subscriptionRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.Error("SubscriptionService.GetById:subscriptionRepo.GetById - Internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	subscription := &domain.Subscription{
		Id:          model.Id,
		ServiceName: model.ServiceName,
		Price:       model.Price,
		UserId:      model.UserId,
		StartDate:   model.StartDate.Format("01-2006"),
	}

	if model.EndDate.Valid {
		subscription.EndDate = new(string)
		*subscription.EndDate = model.EndDate.Time.Format("01-2006")
	}

	s.logger.Info(fmt.Sprintf("Subscription id=%d received successfully", id))

	return subscription, err
}

func (s *SubscriptionService) DeleteById(ctx context.Context, id int) (*domain.Subscription, error) {
	model, err := s.subscriptionRepo.DeleteById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		s.logger.Error("SubscriptionService.DeleteById:subscriptionRepo.DeleteById - Internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	subscription := &domain.Subscription{
		Id:          model.Id,
		ServiceName: model.ServiceName,
		Price:       model.Price,
		UserId:      model.UserId,
		StartDate:   model.StartDate.Format("01-2006"),
	}

	if model.EndDate.Valid {
		subscription.EndDate = new(string)
		*subscription.EndDate = model.EndDate.Time.Format("01-2006")
	}

	s.logger.Info(fmt.Sprintf("Subscription id=%d deleted successfully", id))

	return subscription, err
}

func (s *SubscriptionService) Update(ctx context.Context, data *domain.SubscriptionUpdate) (*domain.Subscription, error) {
	m := &models.SubscriptionUpdate{Id: data.Id}
	if data.ServiceName != nil {
		m.ServiceName = sql.NullString{String: *data.ServiceName, Valid: true}
	}
	if data.Price != nil {
		m.Price = sql.NullInt32{Int32: int32(*data.Price), Valid: true}
	}
	if data.StartDate != nil {
		startDate, _ := time.Parse("01-2006", *data.StartDate)
		m.StartDate = sql.NullTime{Time: startDate, Valid: true}
	}
	if data.EndDate != nil {
		endDate, _ := time.Parse("01-2006", *data.EndDate)
		m.EndDate = sql.NullTime{Time: endDate, Valid: true}
	}

	model, err := s.subscriptionRepo.Update(ctx, m)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		} else if errors.Is(err, repository.ErrIncorrectTime) {
			return nil, ErrIncorrectTime
		}
		s.logger.Error("SubscriptionService.Update:subscriptionRepo.Update - Internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	subscription := &domain.Subscription{
		Id:          model.Id,
		ServiceName: model.ServiceName,
		Price:       model.Price,
		UserId:      model.UserId,
		StartDate:   model.StartDate.Format("01-2006"),
	}

	if model.EndDate.Valid {
		subscription.EndDate = new(string)
		*subscription.EndDate = model.EndDate.Time.Format("01-2006")
	}

	s.logger.Info(fmt.Sprintf("Subscription id=%d update successfully", data.Id))

	return subscription, err
}

func (s *SubscriptionService) GetPriceByFilter(ctx context.Context, user_id *uuid.UUID, service_name *string, startDate, endDate string) (int, error) {
	parsedStart, _ := time.Parse("01-2006", startDate)
	parsedEnd, _ := time.Parse("01-2006", endDate)
	if parsedEnd.Before(parsedStart) {
		return 0, ErrIncorrectTime
	}

	total, err := s.subscriptionRepo.GetPriceByFilter(ctx, user_id, service_name, parsedStart, parsedEnd)
	if err != nil {
		s.logger.Error("SubscriptionService.GetPriceByFilter:subscriptionRepo.GetPriceByFilter - Internal error", slog.String("error", err.Error()))
		return 0, ErrInternal
	}
	total *= int(parsedEnd.Month()) - int(parsedStart.Month()) + 1
	s.logger.Info(fmt.Sprintf("Total cost for the period from %s to %s is %d", startDate, endDate, total))

	return total, err
}
