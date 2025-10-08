package postgres_storage

import (
	"database/sql"

	"github.com/kasparovgs/subscription-aggregation-service/domain"

	sq "github.com/Masterminds/squirrel"
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

func (ps *SubcriptionDB) Close() error {
	if ps.db != nil {
		return ps.db.Close()
	}
	return nil
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

func (ps *SubcriptionDB) GetListOfSubscriptions(filter *domain.SubscriptionFilter) ([]domain.Subscription, error) {
	builder := sq.Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		PlaceholderFormat(sq.Dollar)

	if filter.UserID != nil {
		builder = builder.Where(sq.Eq{"user_id": *filter.UserID})
	}
	if filter.ServiceName != nil {
		builder = builder.Where(sq.Eq{"service_name": *filter.ServiceName})
	}
	if filter.Price != nil {
		builder = builder.Where(sq.Eq{"price": *filter.Price})
	}
	if filter.StartDate != nil {
		builder = builder.Where(sq.GtOrEq{"start_date": *filter.StartDate})
	}
	if filter.EndDate != nil {
		builder = builder.Where(sq.LtOrEq{"end_date": *filter.EndDate})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		if err := rows.Scan(&sub.SubscriptionID, &sub.ServiceName, &sub.Price,
			&sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		result = append(result, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (ps *SubcriptionDB) GetTotalCost(filter *domain.TotalCostFilter) ([]domain.Subscription, error) {
	builder := sq.Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		Where("start_date <= ?", filter.EndDate).
		Where("(end_date IS NULL OR end_date >= ?)", filter.StartDate).
		PlaceholderFormat(sq.Dollar)

	if filter.UserID != nil {
		builder = builder.Where(sq.Eq{"user_id": *filter.UserID})
	}
	if filter.ServiceName != nil {
		builder = builder.Where(sq.Eq{"service_name": *filter.ServiceName})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		err = rows.Scan(&s.SubscriptionID, &s.ServiceName,
			&s.Price, &s.UserID, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
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

func (ps *SubcriptionDB) DeleteSubscriptionByID(subs *domain.Subscription) error {
	if !ps.IsExist(subs.SubscriptionID) {
		return domain.ErrNotFound("subscription not found")
	}
	query := `DELETE FROM subscriptions WHERE id = $1
			  RETURNING service_name, price, user_id, start_date, end_date`
	err := ps.db.QueryRow(query, subs.SubscriptionID).Scan(&subs.ServiceName, &subs.Price, &subs.UserID, &subs.StartDate, &subs.EndDate)
	if err == sql.ErrNoRows {
		return domain.ErrNotFound("subscription not found")
	}
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
