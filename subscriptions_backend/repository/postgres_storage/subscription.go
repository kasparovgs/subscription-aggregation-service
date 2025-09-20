package postgresstorage

import (
	"database/sql"
	"subscriptions_backend/domain"
)

type SubcriptionDB struct {
	db *sql.DB
}

func NewSubscriptionDB(connStr string) (*SubcriptionDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &SubcriptionDB{db: db}, nil
}

func (ps *SubcriptionDB) CreateSubscription(subs *domain.Subscription) error {

}
