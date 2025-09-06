package httpdto

// swagger:model CreateSubscriptionRequest
type CreateSubscriptionRequest struct {
	// Service name (human readable)
	// example: Netflix
	ServiceName string `json:"service_name" binding:"required" example:"Netflix"`

	// Monthly price in whole rubles
	// example: 499
	Price int `json:"price" binding:"required,gte=0" example:"499"`

	// Owner user id (UUIDv4)
	// example: 1c9d4f8b-f0f1-4b9a-8f5e-6e9a0b7f8d12
	UserID string `json:"user_id" binding:"required,uuid" example:"1c9d4f8b-f0f1-4b9a-8f5e-6e9a0b7f8d12"`

	// month-year format: "07-2025"
	// example: 07-2025
	StartDate string `json:"start_date" binding:"required" example:"07-2025"`

	// optional month-year like "12-2025"
	// example: 12-2025
	EndDate *string `json:"end_date,omitempty" example:"12-2025"`
}

// swagger:model UpdateSubscriptionRequest
type UpdateSubscriptionRequest struct {
	// example: Spotify
	ServiceName *string `json:"service_name,omitempty" example:"Spotify"`
	// example: 299
	Price *int `json:"price,omitempty" example:"299"`
	// example: 07-2025
	StartDate *string `json:"start_date,omitempty" example:"07-2025"`
	// To explicitly clear end_date send empty string "".
	// example: 12-2025
	EndDate *string `json:"end_date,omitempty" example:"12-2025"`
}

// swagger:model TotalResponse
type TotalResponse struct {
	// Total amount in whole rubles for the given period (sum over months * price).
	// example: 1497
	Total int64 `json:"total" example:"1497"`
}
