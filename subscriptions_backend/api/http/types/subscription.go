package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

func (r PostCreateSubscriptionRequest) ToDTO() (PostCreateSubscriptionDTO, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return PostCreateSubscriptionDTO{}, err
	}

	start, err := parseMonthYear(r.StartDate)
	if err != nil {
		return PostCreateSubscriptionDTO{}, err
	}

	var end *time.Time
	if r.EndDate != nil {
		parsedEnd, err := parseMonthYear(*r.EndDate)
		if err != nil {
			return PostCreateSubscriptionDTO{}, err
		}
		end = &parsedEnd
	}
	return PostCreateSubscriptionDTO{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func CreatePostSubscriptionRequest(r *http.Request) (*PostCreateSubscriptionRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}

	defer r.Body.Close()

	var req PostCreateSubscriptionRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

func parseMonthYear(s string) (time.Time, error) {
	layout := "01-2006"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
