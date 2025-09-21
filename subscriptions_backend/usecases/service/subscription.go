package service

import (
	"log/slog"
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
		slog.Error("failed to create subscription in repository",
			"error", err,
			"user_id", subs.UserID,
			"service_name", subs.ServiceName,
		)
		return uuid.Nil, err
	}

	slog.Info("subscription created",
		"layer", "service",
		"subscription_id", subscriptionID,
		"user_id", subs.UserID,
		"service_name", subs.ServiceName,
	)
	return subscriptionID, nil
}
