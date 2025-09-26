package service

import (
	"log/slog"
	"time"

	"github.com/kasparovgs/subscription-aggregation-service/domain"

	"github.com/kasparovgs/subscription-aggregation-service/repository"

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

func (s *Subcription) DeleteSubscriptionByID(subs *domain.Subscription) (*domain.Subscription, error) {
	err := s.subscriptionRepo.DeleteSubscriptionByID(subs)
	if err != nil {
		slog.Error("failed to delete subscription from repository",
			"error", err,
			"subscription_id", subs.SubscriptionID,
		)
		return nil, err
	}
	slog.Info("subscription deleted from repo",
		"layer", "service",
		"subscription_id", subs.SubscriptionID,
		"user_id", subs.UserID,
		"service_name", subs.ServiceName)
	return subs, nil
}

func (s *Subcription) GetListOfSubscriptions(filter *domain.SubscriptionFilter) ([]domain.Subscription, error) {
	if filter == nil {
		slog.Error("failed to get list by nil filter")
		return nil, domain.ErrBadRequest("failed to get list by nil filter")
	}
	if filter.StartDate != nil && filter.EndDate != nil &&
		filter.StartDate.After(*filter.EndDate) {
		slog.Error("start date cannot be after end date", "layer", "service")
		return nil, domain.ErrBadRequest("start date cannot be after end date")
	}
	list, err := s.subscriptionRepo.GetListOfSubscriptions(filter)
	if err != nil {
		slog.Error("failed to get list of subscriptions by filter", "layer", "service", "error", err)
		return nil, err
	}
	slog.Info("list of subscriptions by filter successfully found", "layer", "service")
	return list, nil
}

func (s *Subcription) GetTotalCost(filter *domain.TotalCostFilter) (int, error) {
	if filter == nil {
		slog.Error("failed to get total cost by nil filter")
		return 0, domain.ErrBadRequest("failed to get list by nil filter")
	}
	if filter.StartDate.After(filter.EndDate) {
		slog.Error("start date cannot be after end date", "layer", "service")
		return 0, domain.ErrBadRequest("start date cannot be after end date")
	}
	subs, err := s.subscriptionRepo.GetTotalCost(filter)
	if err != nil {
		slog.Error("failed to get total cost of subscriptions by filter", "layer", "service", "error", err)
		return 0, err
	}

	var totalCost int
	for _, sub := range subs {
		totalCost += costForPeriod(&sub, filter.StartDate, filter.EndDate)
	}
	slog.Info("total cost of subscriptions by filter successfully found",
		"layer", "service",
		"total_cost", totalCost)

	return totalCost, nil
}

func costForPeriod(sub *domain.Subscription, periodStart, periodEnd time.Time) int {
	start := maxTime(sub.StartDate, periodStart)

	var end time.Time
	if sub.EndDate != nil {
		end = minTime(*sub.EndDate, periodEnd)
	} else {
		end = periodEnd
	}

	if end.Before(start) {
		return 0
	}

	months := countMonths(start, end)
	return months * sub.Price
}

func countMonths(start, end time.Time) int {
	y1, m1, _ := start.Date()
	y2, m2, _ := end.Date()
	return (y2-y1)*12 + int(m2-m1) + 1
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
