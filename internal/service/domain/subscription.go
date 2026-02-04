package domain

import (
	"github.com/google/uuid"
)

type Subscription struct {
	Id          int
	ServiceName string
	Price       int
	UserId      uuid.UUID
	StartDate   string
	EndDate     *string
}

type SubscriptionCreate struct {
	ServiceName string
	Price       int
	UserId      uuid.UUID
	StartDate   string
	EndDate     *string
}

type SubscriptionUpdate struct {
	Id          int
	ServiceName *string
	Price       *int
	StartDate   *string
	EndDate     *string
}
