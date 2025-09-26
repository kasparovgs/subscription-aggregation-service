package repository

import (
	"github.com/kasparovgs/subscription-aggregation-service/domain"

	"github.com/google/uuid"
)

type SubscriptionDB interface {
	CreateSubscription(subs *domain.Subscription) error
	GetSubscriptionByID(subscriptionID uuid.UUID) (*domain.Subscription, error)
	GetListOfSubscriptions(filter *domain.SubscriptionFilter) ([]domain.Subscription, error)
	GetTotalCost(filter *domain.TotalCostFilter) ([]domain.Subscription, error)
	PatchSubscriptionByID(subs *domain.Subscription) error
	DeleteSubscriptionByID(subs *domain.Subscription) error
	IsExist(subscriptionID uuid.UUID) bool
}
