package postgres_storage

import (
	"database/sql"
	"subscriptions_backend/domain"

	"github.com/google/uuid"
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

func (ps *SubcriptionDB) GetSubscriptionByID(subscriptionID uuid.UUID) (*domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	var subs domain.Subscription
	err := ps.db.QueryRow(query, subscriptionID).Scan(&subs.SubscriptionID,
		&subs.ServiceName,
		&subs.Price,
		&subs.UserID,
		&subs.StartDate,
		&subs.EndDate)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound("subscription not found")
	}
	if err != nil {
		return nil, err
	}

	return &subs, nil
}

func (ps *SubcriptionDB) PatchSubscriptionByID(subs *domain.Subscription) error {
	if !ps.IsExist(subs.SubscriptionID) {
		return domain.ErrNotFound("subscription not found")
	}
	query := `UPDATE subscriptions SET service_name = COALESCE($1, service_name),
         			 price = COALESCE($2, price), end_date = COALESCE($3, end_date)
     				 WHERE id = $4`
	_, err := ps.db.Exec(query, subs.ServiceName, subs.Price, subs.EndDate, subs.SubscriptionID)
	if err != nil {
		return err
	}
	return nil
}

func (ps *SubcriptionDB) IsExist(subscriptionID uuid.UUID) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM subscriptions WHERE id = $1)`
	_ = ps.db.QueryRow(query, subscriptionID).Scan(&exists)
	return exists
}
