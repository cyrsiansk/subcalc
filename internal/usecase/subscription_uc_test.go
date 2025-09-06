package usecase

import (
	"context"
	"errors"
	"subcalc/internal/domain"
	"subcalc/internal/repository"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeRepo struct {
	sumReturn int64
	sumErr    error

	lastFilter repository.SubscriptionFilter
}

func (f *fakeRepo) Count(ctx context.Context, filter repository.SubscriptionFilter) (int64, error) {
	return 0, nil
}
func (f *fakeRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	return nil
}
func (f *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return nil, nil
}
func (f *fakeRepo) Update(ctx context.Context, sub *domain.Subscription) error {
	return nil
}
func (f *fakeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (f *fakeRepo) List(ctx context.Context, filter repository.SubscriptionFilter) ([]*domain.Subscription, error) {
	return nil, nil
}
func (f *fakeRepo) FindForPeriod(ctx context.Context, filter repository.SubscriptionFilter) ([]*domain.Subscription, error) {
	return nil, nil
}
func (f *fakeRepo) SumForPeriod(ctx context.Context, filter repository.SubscriptionFilter) (int64, error) {
	f.lastFilter = filter
	return f.sumReturn, f.sumErr
}

func TestSumSubscriptions_NoPeriod_ReturnsZero(t *testing.T) {
	fr := &fakeRepo{sumReturn: 12345}
	uc := NewSubscriptionUsecase(fr)

	total, err := uc.SumSubscriptions(context.Background(), repository.SubscriptionFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected 0 when no period provided, got %d", total)
	}
}

func TestSumSubscriptions_DelegatesToRepo(t *testing.T) {
	expected := int64(9999)
	fr := &fakeRepo{sumReturn: expected}
	uc := NewSubscriptionUsecase(fr)

	from := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)
	uid := uuid.New()
	sname := "Yandex Plus"

	filter := repository.SubscriptionFilter{
		UserID:      &uid,
		ServiceName: &sname,
		From:        &from,
		To:          &to,
	}

	total, err := uc.SumSubscriptions(context.Background(), filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != expected {
		t.Fatalf("expected %d, got %d", expected, total)
	}

	if fr.lastFilter.UserID == nil || *fr.lastFilter.UserID != uid {
		t.Fatalf("user_id not passed to repo correctly")
	}
	if fr.lastFilter.ServiceName == nil || *fr.lastFilter.ServiceName != sname {
		t.Fatalf("service_name not passed to repo correctly")
	}
	if fr.lastFilter.From == nil || !fr.lastFilter.From.Equal(from) {
		t.Fatalf("from not passed correctly")
	}
	if fr.lastFilter.To == nil || !fr.lastFilter.To.Equal(to) {
		t.Fatalf("to not passed correctly")
	}
}

func TestSumSubscriptions_RepoError_Propagates(t *testing.T) {
	fr := &fakeRepo{sumErr: errors.New("db failing")}
	uc := NewSubscriptionUsecase(fr)

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

	total, err := uc.SumSubscriptions(context.Background(), repository.SubscriptionFilter{From: &from, To: &to})
	if err == nil {
		t.Fatalf("expected error from repository, got nil")
	}
	if total != 0 {
		t.Fatalf("expected total 0 on error, got %d", total)
	}
}
