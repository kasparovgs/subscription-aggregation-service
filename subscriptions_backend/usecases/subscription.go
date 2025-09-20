package usecases

import (
	"subscriptions_backend/domain"

	"github.com/google/uuid"
)

type Subcription interface {
	CreateSubscription(subs *domain.Subscription) (uuid.UUID, error)
}
