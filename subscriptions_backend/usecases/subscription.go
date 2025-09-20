package usecases

import (
	"time"

	"github.com/google/uuid"
)

type Subcription interface {
	CreateSubscription(servName string,
		price int,
		userID uuid.UUID,
		startDate time.Time,
		endDate *time.Time) (uuid.UUID, error)
}
