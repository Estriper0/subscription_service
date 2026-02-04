package dto

import (
	"github.com/google/uuid"
)

// Subscription модель подписки
type Subscription struct {
	Id          int       `json:"id" example:"1"`
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       int       `json:"price" example:"599"`
	UserId      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   string    `json:"start_date" example:"01-2026"`
	EndDate     *string   `json:"end_date" example:"05-2026"`
}

// SubscriptionCreateRequest запрос на создание подписки
type SubscriptionCreateRequest struct {
	ServiceName string  `json:"service_name" validate:"required,lte=100" example:"Netflix"`
	Price       int     `json:"price" validate:"required,gte=0" example:"599"`
	UserId      string  `json:"user_id" validate:"required,uuid4" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   string  `json:"start_date" validate:"required,date" example:"01-2026"`
	EndDate     *string `json:"end_date" validate:"omitempty,date" example:"05-2026"`
}

// SubscriptionUpdateRequest запрос на обновление подписки
type SubscriptionUpdateRequest struct {
	Id          int     `json:"id" validate:"required,gte=0" example:"1"`
	ServiceName *string `json:"service_name" validate:"omitempty,lte=100" example:"Netflix Premium"`
	Price       *int    `json:"price" validate:"omitempty,gte=0" example:"699"`
	StartDate   *string `json:"start_date" validate:"omitempty,date" example:"01-2026"`
	EndDate     *string `json:"end_date" validate:"omitempty,date" example:"05-2026"`
}

// SubscriptionFilterRequest запрос для фильтрации подписок
type SubscriptionFilterRequest struct {
	StartDate string `json:"start_date" validate:"required,date" example:"01-2026"`
	EndDate   string `json:"end_date" validate:"required,date" example:"05-2026"`
}
