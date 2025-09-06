package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// swagger:model Subscription
type Subscription struct {
	// Unique subscription id (UUIDv4)
	// example: 3fa85f64-5717-4562-b3fc-2c963f66afa6
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`

	// Human-friendly service name
	// example: Netflix
	ServiceName string `json:"service_name" example:"Netflix"`

	// Monthly price in whole rubles (integer)
	// example: 499
	Price int `json:"price" example:"499"`

	// Owner user id (UUIDv4)
	// example: 1c9d4f8b-f0f1-4b9a-8f5e-6e9a0b7f8d12
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;index" example:"1c9d4f8b-f0f1-4b9a-8f5e-6e9a0b7f8d12"`

	// Start date (month precision). Rendered в JSON как "MM-YYYY".
	// example: 07-2025
	// swagger type: string
	StartDate time.Time `json:"start_date" swaggertype:"string" example:"07-2025"`

	// Optional end date (month precision). Rendered в JSON как "MM-YYYY".
	// example: 12-2025
	// swagger type: string
	EndDate *time.Time `json:"end_date,omitempty" swaggertype:"string" example:"12-2025"`

	// Timestamps (server side). RFC3339.
	// example: 2025-07-01T12:00:00Z
	CreatedAt time.Time `json:"created_at" example:"2025-07-01T12:00:00Z"`

	// example: 2025-07-01T12:00:00Z
	UpdatedAt time.Time `json:"updated_at" example:"2025-07-01T12:00:00Z"`
}

func (s Subscription) MarshalJSON() ([]byte, error) {
	type aux struct {
		ID          uuid.UUID `json:"id"`
		ServiceName string    `json:"service_name"`
		Price       int       `json:"price"`
		UserID      uuid.UUID `json:"user_id"`
		StartDate   string    `json:"start_date"`
		EndDate     *string   `json:"end_date,omitempty"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	start := fmt.Sprintf("%02d-%04d", s.StartDate.Month(), s.StartDate.Year())
	var end *string
	if s.EndDate != nil {
		t := fmt.Sprintf("%02d-%04d", s.EndDate.Month(), s.EndDate.Year())
		end = &t
	}

	a := aux{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   start,
		EndDate:     end,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}

	return json.Marshal(a)
}
