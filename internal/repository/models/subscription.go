package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	Id          int
	ServiceName string
	Price       int
	UserId      uuid.UUID
	StartDate   time.Time
	EndDate     sql.NullTime
}

type SubscriptionCreate struct {
	ServiceName string
	Price       int
	UserId      uuid.UUID
	StartDate   time.Time
	EndDate     sql.NullTime
}

type SubscriptionUpdate struct {
	Id          int
	ServiceName sql.NullString
	Price       sql.NullInt32
	StartDate   sql.NullTime
	EndDate     sql.NullTime
}
