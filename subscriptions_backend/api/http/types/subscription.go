package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"subscriptions_backend/domain"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PostCreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type PostCreateSubscriptionDTO struct {
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

func (r PostCreateSubscriptionRequest) ToDomain() (*domain.Subscription, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding uuid: %v", err))
	}

	start, err := parseMonthYear(r.StartDate)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding startDate: %v", err))
	}

	var end *time.Time
	if r.EndDate != nil {
		parsedEnd, err := parseMonthYear(*r.EndDate)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding endDate: %v", err))
		}
		end = &parsedEnd
	}
	return &domain.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func CreatePostSubscriptionHandlerRequest(r *http.Request) (*PostCreateSubscriptionRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding json: %v", err))
	}

	defer r.Body.Close()

	var req PostCreateSubscriptionRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding json: %v", err))
	}
	return &req, nil
}

type PostCreateSubscriptionResponse struct {
	SubscriptionID uuid.UUID `json:"subscription_id"`
}

func GetSubscriptionByIDHandlerRequest(r *http.Request) (*domain.Subscription, error) {
	subIDStr := chi.URLParam(r, "subscription_id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding uuid: %v", err))
	}
	subs := domain.Subscription{SubscriptionID: subID}
	return &subs, nil
}

type GetSubscriptionByIDResponse struct {
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	ServiceName    string     `json:"service_name"`
	Price          int        `json:"price"`
	UserID         uuid.UUID  `json:"user_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}

func parseMonthYear(s string) (time.Time, error) {
	layout := "01-2006"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
