package repository

import (
	"context"
	"subcalc/internal/domain"
	"time"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter SubscriptionFilter) ([]*domain.Subscription, error)

	FindForPeriod(ctx context.Context, filter SubscriptionFilter) ([]*domain.Subscription, error)
	SumForPeriod(ctx context.Context, filter SubscriptionFilter) (int64, error)
	Count(ctx context.Context, filter SubscriptionFilter) (int64, error)
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	From        *time.Time
	To          *time.Time
	Limit       int
	Offset      int
}
