package postgres_storage

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
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := ps.db.Exec(query, subs.SubscriptionID, subs.ServiceName, subs.Price, subs.UserID, subs.StartDate, subs.EndDate)
	if err != nil {
		return err
	}
	return nil
}
