package repository

import (
	"subscriptions_backend/domain"

	"github.com/google/uuid"
)

type SubscriptionDB interface {
	CreateSubscription(subs *domain.Subscription) error
	GetSubscriptionByID(subscriptionID uuid.UUID) (*domain.Subscription, error)
}
