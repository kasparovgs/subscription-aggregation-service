package repository

import "subscriptions_backend/domain"

type SubscriptionDB interface {
	CreateSubscription(subs *domain.Subscription) error
}
