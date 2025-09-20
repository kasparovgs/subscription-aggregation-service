package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	ServiceName    string     `json:"service_name"`
	Price          int        `json:"price"`
	UserID         uuid.UUID  `json:"user_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
