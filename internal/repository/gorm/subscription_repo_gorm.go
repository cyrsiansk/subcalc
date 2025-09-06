package gormrepo

import (
	"context"
	"errors"
	"subcalc/internal/domain"
	"subcalc/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormSubscription struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	ServiceName string     `json:"service_name" gorm:"type:text;not null"`
	Price       int        `json:"price" gorm:"type:int;not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;index;not null"`
	StartDate   time.Time  `json:"start_date" gorm:"type:date;not null"`
	EndDate     *time.Time `json:"end_date" gorm:"type:date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (g *GormSubscription) TableName() string {
	return "subscriptions"
}

func (g *GormSubscription) ToDomain() *domain.Subscription {
	return &domain.Subscription{
		ID:          g.ID,
		ServiceName: g.ServiceName,
		Price:       g.Price,
		UserID:      g.UserID,
		StartDate:   g.StartDate,
		EndDate:     g.EndDate,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

func FromDomain(d *domain.Subscription) *GormSubscription {
	return &GormSubscription{
		ID:          d.ID,
		ServiceName: d.ServiceName,
		Price:       d.Price,
		UserID:      d.UserID,
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

type repo struct {
	db *gorm.DB
}

func NewGormSubscriptionRepo(db *gorm.DB) repository.SubscriptionRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, sub *domain.Subscription) error {
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}
	g := FromDomain(sub)
	if err := r.db.WithContext(ctx).Create(g).Error; err != nil {
		return err
	}
	sub.ID = g.ID
	sub.CreatedAt = g.CreatedAt
	sub.UpdatedAt = g.UpdatedAt
	return nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var g GormSubscription
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return g.ToDomain(), nil
}

func (r *repo) Update(ctx context.Context, sub *domain.Subscription) error {
	now := time.Now().UTC()
	updates := map[string]interface{}{
		"service_name": sub.ServiceName,
		"price":        sub.Price,
		"user_id":      sub.UserID,
		"start_date":   sub.StartDate,
		"end_date":     sub.EndDate,
		"updated_at":   now,
	}
	if err := r.db.WithContext(ctx).Model(&GormSubscription{}).Where("id = ?", sub.ID).Updates(updates).Error; err != nil {
		return err
	}
	sub.UpdatedAt = now
	return nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&GormSubscription{}, "id = ?", id).Error
}

func (r *repo) List(ctx context.Context, filter repository.SubscriptionFilter) ([]*domain.Subscription, error) {
	var gs []GormSubscription
	q := r.db.WithContext(ctx).Model(&GormSubscription{})

	if filter.ServiceName != nil {
		q = q.Where("service_name = ?", *filter.ServiceName)
	}
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}

	if filter.From != nil && filter.To != nil {
		q = q.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", *filter.To, *filter.From)
	} else if filter.From != nil {
		q = q.Where("end_date IS NULL OR end_date >= ?", *filter.From)
	} else if filter.To != nil {
		q = q.Where("start_date <= ?", *filter.To)
	}

	if filter.Limit == 0 {
		filter.Limit = 100
	}
	if filter.Offset > 0 {
		q = q.Offset(filter.Offset)
	}
	if err := q.Limit(filter.Limit).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.Subscription, 0, len(gs))
	for _, g := range gs {
		out = append(out, g.ToDomain())
	}
	return out, nil
}

func (r *repo) FindForPeriod(ctx context.Context, filter repository.SubscriptionFilter) ([]*domain.Subscription, error) {
	var gs []GormSubscription
	q := r.db.WithContext(ctx).Model(&GormSubscription{})
	if filter.ServiceName != nil {
		q = q.Where("service_name = ?", *filter.ServiceName)
	}
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.From != nil && filter.To != nil {
		q = q.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", *filter.To, *filter.From)
	} else if filter.From != nil {
		q = q.Where("end_date IS NULL OR end_date >= ?", *filter.From)
	} else if filter.To != nil {
		q = q.Where("start_date <= ?", *filter.To)
	}
	if filter.Limit == 0 {
		filter.Limit = 1000
	}
	if filter.Offset > 0 {
		q = q.Offset(filter.Offset)
	}
	if err := q.Limit(filter.Limit).Find(&gs).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.Subscription, 0, len(gs))
	for _, g := range gs {
		out = append(out, g.ToDomain())
	}
	return out, nil
}

func (r *repo) Count(ctx context.Context, filter repository.SubscriptionFilter) (int64, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&GormSubscription{})
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.ServiceName != nil {
		q = q.Where("service_name = ?", *filter.ServiceName)
	}
	if filter.From != nil && filter.To != nil {
		q = q.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", *filter.To, *filter.From)
	} else if filter.From != nil {
		q = q.Where("end_date IS NULL OR end_date >= ?", *filter.From)
	} else if filter.To != nil {
		q = q.Where("start_date <= ?", *filter.To)
	}
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repo) SumForPeriod(ctx context.Context, filter repository.SubscriptionFilter) (int64, error) {
	if filter.From == nil || filter.To == nil {
		return 0, nil
	}

	from := dateTruncMonth(*filter.From)
	to := dateTruncMonth(*filter.To)

	where := " WHERE start_date <= ?::date AND (end_date IS NULL OR end_date >= ?::date)"
	whereArgs := []interface{}{to, from}

	if filter.ServiceName != nil {
		where = where + " AND service_name = ?"
		whereArgs = append(whereArgs, *filter.ServiceName)
	}
	if filter.UserID != nil {
		where = where + " AND user_id = ?"
		whereArgs = append(whereArgs, *filter.UserID)
	}

	base := `
WITH periods AS (
  SELECT
    GREATEST(start_date, ?::date) AS s,
    LEAST(COALESCE(end_date, ?::date), ?::date) AS e,
    price
  FROM subscriptions
  ` + where + `
)
SELECT COALESCE(SUM(price * (
  (DATE_PART('year', AGE(e, s)) * 12) + DATE_PART('month', AGE(e, s)) + 1
)), 0) AS total
FROM periods
WHERE e >= s
`
	args := make([]interface{}, 0, 8)
	args = append(args, from, to, to)
	args = append(args, whereArgs...)

	var res struct {
		Total int64 `gorm:"column:total"`
	}
	if err := r.db.WithContext(ctx).Raw(base, args...).Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Total, nil
}

func dateTruncMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}
