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

func (s *Subcription) GetSubscriptionByID(subscriptionID uuid.UUID) (*domain.Subscription, error) {
	subs, err := s.subscriptionRepo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		slog.Error("failed to get subscription from repository",
			"error", err,
			"subscription_id", subscriptionID,
		)
		return nil, err
	}

	slog.Info("subscription received from repo",
		"layer", "service",
		"subscription_id", subscriptionID,
		"user_id", subs.UserID,
		"service_name", subs.ServiceName,
	)
	return subs, nil
}

func (s *Subcription) PatchSubscriptionByID(subs *domain.Subscription) (*domain.Subscription, error) {
	err := s.subscriptionRepo.PatchSubscriptionByID(subs)
	if err != nil {
		slog.Error("failed to patch subscription in repository",
			"error", err,
			"subscription_id", subs.SubscriptionID,
		)
		return nil, err
	}

	subs, _ = s.subscriptionRepo.GetSubscriptionByID(subs.SubscriptionID)
	slog.Info("subscription patched in repo",
		"layer", "service",
		"subscription_id", subs.SubscriptionID,
		"user_id", subs.UserID,
		"service_name", subs.ServiceName)
	return subs, nil
}
