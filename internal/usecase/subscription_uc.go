package usecase

import (
	"context"
	"subcalc/internal/domain"
	"subcalc/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionUsecase interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter repository.SubscriptionFilter) ([]*domain.Subscription, error)
	SumSubscriptions(ctx context.Context, filter repository.SubscriptionFilter) (int64, error)
	Count(ctx context.Context, filter repository.SubscriptionFilter) (int64, error)
}

type subscriptionUC struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionUsecase(repo repository.SubscriptionRepository) SubscriptionUsecase {
	return &subscriptionUC{repo: repo}
}

func (u *subscriptionUC) Create(ctx context.Context, sub *domain.Subscription) error {
	return u.repo.Create(ctx, sub)
}

func (u *subscriptionUC) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *subscriptionUC) Update(ctx context.Context, sub *domain.Subscription) error {
	return u.repo.Update(ctx, sub)
}

func (u *subscriptionUC) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

func (u *subscriptionUC) List(ctx context.Context, filter repository.SubscriptionFilter) ([]*domain.Subscription, error) {
	return u.repo.List(ctx, filter)
}

func (u *subscriptionUC) SumSubscriptions(ctx context.Context, filter repository.SubscriptionFilter) (int64, error) {
	if filter.From == nil || filter.To == nil {
		return 0, nil
	}
	return u.repo.SumForPeriod(ctx, filter)
}

func (u *subscriptionUC) Count(ctx context.Context, filter repository.SubscriptionFilter) (int64, error) {
	return u.repo.Count(ctx, filter)
}
