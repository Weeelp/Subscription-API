package models

type Subscription struct {
	ID          int     `db:"id" json:"id"`
	ServiceName string  `db:"service_name" json:"service_name"`
	Price       int     `db:"price" json:"price"`
	UserID      string  `db:"user_id" json:"user_id"`
	StartDate   string  `db:"start_date" json:"start_date"`
	EndDate     *string `db:"end_date" json:"end_date"`
}
