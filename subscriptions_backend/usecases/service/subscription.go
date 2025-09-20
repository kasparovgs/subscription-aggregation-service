package service

import (
	"subscriptions_backend/domain"
	"subscriptions_backend/repository"

	"github.com/google/uuid"
)

type Subcription struct {
	subscriptionRepo repository.SubscriptionDB
}

func NewSubscription(subsRepo repository.SubscriptionDB) *Subcription {
	return &Subcription{subscriptionRepo: subsRepo}
}

func (s *Subcription) CreateSubscription(subs *domain.Subscription) (uuid.UUID, error) {
	subscriptionID := uuid.New()
	subs.SubscriptionID = subscriptionID
	err := s.subscriptionRepo.CreateSubscription(subs)
	if err != nil {
		return uuid.Nil, err
	}
	return subscriptionID, nil
}
