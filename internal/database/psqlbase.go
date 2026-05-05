package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"pupupu/internal/config"
	"pupupu/internal/logger"
)

type PSQL struct {
	conn *sql.DB
}

type Subscription struct {
	ID          int     `db:"id" json:"id"`
	ServiceName string  `db:"service_name" json:"service_name"`
	Price       int     `db:"price" json:"price"`
	UserID      string  `db:"user_id" json:"user_id"`
	StartDate   string  `db:"start_date" json:"start_date"`
	EndDate     *string `db:"end_date" json:"end_date"`
}

func (p *PSQL) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *PSQL) GetConn() *sql.DB {
	return p.conn
}

func Init(cfg *config.Config) *PSQL {
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
	return &PSQL{
		conn: conn,
	}
}

func (p *PSQL) CreateSub(s Subscription) (int, error) {
	var id int
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := p.conn.QueryRow(query, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert subscription: %w", err)
	}
	return id, nil
}

func (p *PSQL) GetSubByID(id int) (Subscription, error) {
	var s Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date 
			FROM subscriptions WHERE id = $1`

	err := p.conn.QueryRow(query, id).Scan(
		&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return Subscription{}, fmt.Errorf("subscription not found")
		}
		return Subscription{}, fmt.Errorf("query error: %w", err)
	}

	return s, nil
}

func (p *PSQL) GetAllSubs() ([]Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`

	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var subs []Subscription

	for rows.Next() {
		var s Subscription
		err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (p *PSQL) GetTotal(userID, serviceName, period string) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions 
            WHERE user_id = $1 AND service_name = $2 AND start_date = $3`
	err := p.conn.QueryRow(query, userID, serviceName, period).Scan(&total)
	return total, err
}

func (p *PSQL) UpdateSub(id int, s Subscription) error {
	query := `UPDATE subscriptions 
			SET service_name = $1, price = $2, start_date = $3, end_date = $4 
			WHERE id = $5`

	res, err := p.conn.Exec(query, s.ServiceName, s.Price, s.StartDate, s.EndDate, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

func (p *PSQL) DeleteSub(id int) error {
	query := `DELETE FROM subscriptions
			WHERE id = $1`

	res, err := p.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("subscription with id %d not found", id)
	}

	return nil
}
