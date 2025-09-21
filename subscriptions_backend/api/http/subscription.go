package http

import (
	"log/slog"
	"net/http"
	"subscriptions_backend/api/http/types"
	"subscriptions_backend/usecases"

	"github.com/go-chi/chi/v5"
)

// Subscription represents an HTTP handler for managing subscriptions.
type Subscription struct {
	service usecases.Subcription
}

// NewHandler creates a new instance of Subscription.
func NewSubscriptionHandler(service usecases.Subcription) *Subscription {
	return &Subscription{service: service}
}

// @Summary Create a new subscription
// @Description Create a new subscription and issue their subscriptionID
// @Tags subscription
// @Accept  json
// @Produce json
// @Param request body types.PostCreateSubscriptionRequest true "login and password"
// @Success 201 {string} types.PostCreateSubscriptionResponse
// @Failure 400 {string} string "Bad request"
// @Router /subscriptions [post]
func (s *Subscription) postCreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostSubscriptionHandlerRequest(r)
	if err != nil {
		slog.Warn("failed to parse request", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	subscription, err := req.ToDomain()
	if err != nil {
		slog.Warn("failed to convert request to domain", "error", err)
		types.ProcessError(w, err, nil)
		return
	}

	subID, err := s.service.CreateSubscription(subscription)
	if err != nil {
		slog.Error("failed to create subscription in service", "error", err)
		types.ProcessError(w, err, nil)
		return
	}

	slog.Info("subscription created", "subscription_id", subID)
	types.ProcessError(w, err, &types.PostCreateSubscriptionResponse{SubscriptionID: subID})
}

// @Summary Get a subscription
// @Description Get a subscription by their subscriptionID
// @Tags subscription
// @Accept  json
// @Produce json
// @Param subscription_id path string true "UUID of the subscription" format(uuid)
// @Success 200 {string} types.GetSubscriptionByIDResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Subscription not found"
// @Router /subscriptions/{subscription_id} [get]
func (s *Subscription) getSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	subs, err := types.GetSubscriptionByIDHandlerRequest(r)
	if err != nil {
		slog.Warn("failed to parse request", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	subs, err = s.service.GetSubscriptionByID(subs.SubscriptionID)
	if err != nil {
		slog.Error("failed to get subscription by subscriptionID", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	slog.Info("subscription received", "subscription_id", subs.SubscriptionID)
	types.ProcessError(w, err, &types.GetSubscriptionByIDResponse{SubscriptionID: subs.SubscriptionID,
		ServiceName: subs.ServiceName,
		Price:       subs.Price,
		UserID:      subs.UserID,
		StartDate:   subs.StartDate,
		EndDate:     subs.EndDate,
	})
}

// @Summary Patch a subscription
// @Description Patch a subscription by their subscriptionID
// @Tags subscription
// @Accept  json
// @Produce json
// @Param subscription_id path string true "UUID of the subscription" format(uuid)
// @Param request body types.PatchSubscriptionByIDRequest true "Fields to update"
// @Success 200 {string} types.GetSubscriptionByIDResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Subscription not found"
// @Router /subscriptions/{subscription_id} [patch]
func (s *Subscription) patchSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.PatchSubscriptionByIDHandlerRequest(r)
	if err != nil {
		slog.Warn("failed to parse request", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	subscription, err := req.ToDomain()
	if err != nil {
		slog.Warn("failed to convert request to domain", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	subs, err := s.service.PatchSubscriptionByID(subscription)
	if err != nil {
		slog.Error("failed to patch subscription by subscriptionID", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	slog.Info("subscription patched", "subscription_id", subscription.SubscriptionID)
	types.ProcessError(w, err, &types.PatchSubscriptionByIDResponse{SubscriptionID: subs.SubscriptionID,
		ServiceName: subs.ServiceName, Price: subs.Price, UserID: subs.UserID, StartDate: subs.StartDate,
		EndDate: subs.EndDate})
}

func (s *Subscription) WithSubscriptionHandlers(r chi.Router) {
	r.Post("/subscriptions", s.postCreateSubscriptionHandler)
	r.Get("/subscriptions/{subscription_id}", s.getSubscriptionByIDHandler)
	r.Patch("/subscriptions/{subscription_id}", s.patchSubscriptionByIDHandler)
}
