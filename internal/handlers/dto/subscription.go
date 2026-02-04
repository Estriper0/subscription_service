package dto

import (
	"github.com/google/uuid"
)

type Subscription struct {
	Id          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
}

type SubscriptionCreateRequest struct {
	ServiceName string  `json:"service_name" validate:"required,lte=100"`
	Price       int     `json:"price" validate:"required,gte=0"`
	UserId      string  `json:"user_id" validate:"required,uuid4"`
	StartDate   string  `json:"start_date" validate:"required,date"`
	EndDate     *string `json:"end_date" validate:"omitempty,date"`
}

type SubscriptionUpdateRequest struct {
	Id          int     `json:"id" validate:"required,gte=0"`
	ServiceName *string `json:"service_name" validate:"omitempty,lte=100"`
	Price       *int    `json:"price" validate:"omitempty,gte=0"`
	StartDate   *string `json:"start_date" validate:"omitempty,date"`
	EndDate     *string `json:"end_date" validate:"omitempty,date"`
}

type SubscriptionFilterRequest struct {
	StartDate string `json:"start_date" validate:"required,date"`
	EndDate   string `json:"end_date" validate:"required,date"`
}
