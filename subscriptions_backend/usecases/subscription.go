package usecases

import (
	"subscriptions_backend/domain"

	"github.com/google/uuid"
)

type Subcription interface {
	CreateSubscription(subs *domain.Subscription) (uuid.UUID, error)
	GetSubscriptionByID(subscriptionID uuid.UUID) (*domain.Subscription, error)
	PatchSubscriptionByID(subs *domain.Subscription) (*domain.Subscription, error)
}
