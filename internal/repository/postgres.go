package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"pupupu/internal/config"
	"pupupu/internal/logger"
	"pupupu/internal/models"
)

var _ SubscriptionRepository = (*PostgresRepo)(nil)

type PostgresRepo struct {
	conn *sql.DB
}

func (p *PostgresRepo) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *PostgresRepo) GetConn() *sql.DB {
	return p.conn
}

func Init(cfg *config.Config) SubscriptionRepository {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PSQL.Host, cfg.PSQL.Port, cfg.PSQL.User, cfg.PSQL.Password, cfg.PSQL.Name, cfg.PSQL.SSLmode)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Ошибка открытия PSQL: ", "error", err)
		return nil
	}

	if err := conn.Ping(); err != nil {
		logger.Log.Error("Ошибка соединения PSQL: ", "error", err)
		return nil
	}

	logger.Log.Info("PSQL connection success")
	return &PostgresRepo{
		conn: conn,
	}
}

func (p *PostgresRepo) CreateSub(ctx context.Context, s models.Subscription) (int, error) {
	var id int
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := p.conn.QueryRowContext(ctx, query, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert: %w", err)
	}
	return id, nil
}

func (p *PostgresRepo) GetSubByID(ctx context.Context, id int) (models.Subscription, error) {
	var s models.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscriptions WHERE id = $1`

	err := p.conn.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Subscription{}, fmt.Errorf("subscription not found")
		}
		return models.Subscription{}, fmt.Errorf("query error: %w", err)
	}

	return s, nil
}

func (p *PostgresRepo) GetAllSubs(ctx context.Context) ([]models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`

	rows, err := p.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var subs []models.Subscription

	for rows.Next() {
		var s models.Subscription
		err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		subs = append(subs, s)
	}

	return subs, rows.Err()
}

func (p *PostgresRepo) GetTotal(ctx context.Context, userID, serviceName, period string) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions 
            WHERE user_id = $1 AND service_name = $2 AND start_date = $3`

	err := p.conn.QueryRowContext(ctx, query, userID, serviceName, period).Scan(&total)
	return total, err
}

func (p *PostgresRepo) UpdateSub(ctx context.Context, id int, s models.Subscription) error {
	query := `UPDATE subscriptions 
			SET service_name = $1, price = $2, start_date = $3, end_date = $4 
			WHERE id = $5`

	res, err := p.conn.ExecContext(ctx, query, s.ServiceName, s.Price, s.StartDate, s.EndDate, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

func (p *PostgresRepo) DeleteSub(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	res, err := p.conn.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}
