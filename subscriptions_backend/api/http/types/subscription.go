package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"subscriptions_backend/domain"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ***** [POST] CreateSubscription *****
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

func (r *PostCreateSubscriptionRequest) ToDomain() (*domain.Subscription, error) {
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

// *************************************

// ***** [GET] GetSubscriptionByID *****

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

// *************************************

// ***** [PATCH] PatchSubscriptionByID *****

type PatchSubscriptionByIDRequest struct {
	SubscriptionID uuid.UUID `json:"-"`
	ServiceName    *string   `json:"service_name,omitempty"`
	Price          *int      `json:"price"`
	EndDate        *string   `json:"end_date,omitempty"`
}

func PatchSubscriptionByIDHandlerRequest(r *http.Request) (*PatchSubscriptionByIDRequest, error) {
	subIDStr := chi.URLParam(r, "subscription_id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding uuid: %v", err))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding json: %v", err))
	}

	defer r.Body.Close()

	var req PatchSubscriptionByIDRequest

	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding json: %v", err))
	}
	if req.ServiceName == nil && req.Price == nil && req.EndDate == nil {
		return nil, domain.ErrBadRequest("no fields to update")
	}
	return &PatchSubscriptionByIDRequest{SubscriptionID: subID, ServiceName: req.ServiceName, Price: req.Price, EndDate: req.EndDate}, nil
}

func (r *PatchSubscriptionByIDRequest) ToDomain() (*domain.Subscription, error) {
	var end *time.Time
	if r.EndDate != nil {
		parsedEnd, err := parseMonthYear(*r.EndDate)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding endDate: %v", err))
		}
		end = &parsedEnd
	}

	subs := &domain.Subscription{SubscriptionID: r.SubscriptionID}
	if r.ServiceName != nil {
		subs.ServiceName = *r.ServiceName
	}
	if r.Price != nil {
		subs.Price = *r.Price
	}
	if end != nil {
		subs.EndDate = end
	}
	return subs, nil
}

type PatchSubscriptionByIDResponse struct {
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	ServiceName    string     `json:"service_name"`
	Price          int        `json:"price"`
	UserID         uuid.UUID  `json:"user_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}

// *****************************************

// ***** [DELETE] DeleteSubscriptionByID *****
func DeleteSubscriptionByIDHandlerRequest(r *http.Request) (*domain.Subscription, error) {
	subIDStr := chi.URLParam(r, "subscription_id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding uuid: %v", err))
	}
	subs := domain.Subscription{SubscriptionID: subID}
	return &subs, nil
}

type DeleteSubscriptionByIDResponse struct {
	SubscriptionID uuid.UUID  `json:"subscription_id"`
	ServiceName    string     `json:"service_name"`
	Price          int        `json:"price"`
	UserID         uuid.UUID  `json:"user_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
}

// *******************************************

// ***** [GET] GetListOfSubscriptions *****

func GetListOfSubscriptionsHandlerRequest(r *http.Request) (*domain.SubscriptionFilter, error) {
	q := r.URL.Query()
	filter := domain.SubscriptionFilter{}

	if u := q.Get("user_id"); u != "" {
		parsedUUID, err := uuid.Parse(u)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding uuid: %v", err))
		}
		filter.UserID = &parsedUUID
	}

	if s := q.Get("start_date"); s != "" {
		parsedStart, err := parseMonthYear(s)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding startDate: %v", err))
		}
		filter.StartDate = &parsedStart
	}
	if e := q.Get("end_date"); e != "" {
		parsedEnd, err := parseMonthYear(e)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding endDate: %v", err))
		}
		filter.EndDate = &parsedEnd
	}

	if p := q.Get("price"); p != "" {
		parsedPrice, err := strconv.Atoi(p)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding price: %v", err))
		}
		filter.Price = &parsedPrice
	}
	if s := q.Get("service_name"); s != "" {
		filter.ServiceName = &s
	}

	return &filter, nil
}

type GetListOfSubscriptionsResponse struct {
	Subscriptions []domain.Subscription `json:"subscriptions"`
}

// ****************************************

// ***** [GET] GetTotalCost *****

func GetTotalCostHandlerRequest(r *http.Request) (*domain.TotalCostFilter, error) {
	q := r.URL.Query()
	req := domain.TotalCostFilter{}

	if u := q.Get("user_id"); u != "" {
		parsedUUID, err := uuid.Parse(u)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding uuid: %v", err))
		}
		req.UserID = &parsedUUID
	}

	if s := q.Get("service_name"); s != "" {
		req.ServiceName = &s
	}

	if s := q.Get("start_date"); s != "" {
		parsedStart, err := parseMonthYear(s)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding startDate: %v", err))
		}
		req.StartDate = parsedStart
	} else {
		return nil, domain.ErrBadRequest("start_date is required for the request")
	}
	if e := q.Get("end_date"); e != "" {
		parsedEnd, err := parseMonthYear(e)
		if err != nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("error while decoding endDate: %v", err))
		}
		req.EndDate = parsedEnd
	} else {
		return nil, domain.ErrBadRequest("start_date is required for the request")
	}
	return &req, nil
}

type GetTotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

// ******************************

func parseMonthYear(s string) (time.Time, error) {
	layout := "01-2006"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
