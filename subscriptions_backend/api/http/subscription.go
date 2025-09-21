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
// @Tags user
// @Accept  json
// @Produce json
// @Param request body types.PostCreateSubscriptionRequest true "login and password"
// @Success 201 {string} types.PostCreateSubscriptionResponse
// @Failure 400 {string} string "Bad request"
// @Router /subscriptions [post]
func (s *Subscription) postCreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostSubscriptionRequest(r)
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

func (s *Subscription) getSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Subscription) WithSubscriptionHandlers(r chi.Router) {
	r.Post("/subscriptions", s.postCreateSubscriptionHandler)
	r.Get("/subscriptions/{subscription_id}", s.getSubscriptionHandler)
}
