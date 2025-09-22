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
// @Success 201 {object} types.PostCreateSubscriptionResponse
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
// @Success 200 {object} types.GetSubscriptionByIDResponse
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
// @Success 200 {object} types.GetSubscriptionByIDResponse
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

// @Summary Delete a subscription
// @Description Delete a subscription by their subscriptionID
// @Tags subscription
// @Accept  json
// @Produce json
// @Param subscription_id path string true "UUID of the subscription" format(uuid)
// @Success 200 {object} types.GetSubscriptionByIDResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Subscription not found"
// @Router /subscriptions/{subscription_id} [delete]
func (s *Subscription) deleteSubscriptionByIDHandler(w http.ResponseWriter, r *http.Request) {
	subs, err := types.GetSubscriptionByIDHandlerRequest(r)
	if err != nil {
		slog.Warn("failed to parse request", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	subs, err = s.service.DeleteSubscriptionByID(subs)
	if err != nil {
		slog.Error("failed to delete subscription by subscriptionID", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	slog.Info("subscription deleted", "subscription_id", subs.SubscriptionID)
	types.ProcessError(w, err, &types.DeleteSubscriptionByIDResponse{SubscriptionID: subs.SubscriptionID,
		ServiceName: subs.ServiceName, Price: subs.Price, UserID: subs.UserID, StartDate: subs.StartDate,
		EndDate: subs.EndDate})
}

// @Summary List subscriptions
// @Description Get a list of subscriptions with the ability to filter
// @Tags subscription
// @Accept  json
// @Produce json
// @Param user_id query string false "userUUID"
// @Param service_name query string false "Service name"
// @Param start_date query string false "Start date (MM-YYYY)"
// @Param end_date query string false "End date (MM-YYYY)"
// @Success 200 {array} types.GetSubscriptionByIDResponse
// @Failure 400 {string} string "Bad request"
// @Router /subscriptions [get]
func (s *Subscription) getListOfSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.GetListOfSubscriptionsHandlerRequest(r)
	if err != nil {
		slog.Warn("failed to parse request", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	list, err := s.service.GetListOfSubscriptions(req)
	if err != nil {
		slog.Error("filed to get list of subscriptions by filter", "error", err)
		types.ProcessError(w, err, nil)
		return
	}
	slog.Info("list of subscriptions by filter successfully found")
	types.ProcessError(w, err, &types.GetListOfSubscriptionsResponse{Subscriptions: list})
}

func (s *Subscription) WithSubscriptionHandlers(r chi.Router) {
	r.Post("/subscriptions", s.postCreateSubscriptionHandler)
	r.Get("/subscriptions/{subscription_id}", s.getSubscriptionByIDHandler)
	r.Get("/subscriptions", s.getListOfSubscriptionsHandler)
	r.Patch("/subscriptions/{subscription_id}", s.patchSubscriptionByIDHandler)
	r.Delete("/subscriptions/{subscription_id}", s.deleteSubscriptionByIDHandler)
}
