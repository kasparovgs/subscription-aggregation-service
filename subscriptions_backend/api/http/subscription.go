package http

import (
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

func (s *Subscription) postCreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostSubscriptionRequest(r)
	if err != nil {
		types.ProcessError(w, err, nil)
		return
	}
	dto, err := req.ToDTO()
	if err != nil {
		types.ProcessError(w, err, nil)
		return
	}

	subID, err := s.service.CreateSubscription()
}

func (s *Subscription) WithSubscriptionHandlers(r chi.Router) {
	r.Post("/subscriptions", s.postCreateSubscriptionHandler)
}
