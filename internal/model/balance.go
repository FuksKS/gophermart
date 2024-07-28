package model

import "time"

type Balance struct {
	Current   float64 `json:"current" db:"current_balance"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
}

type Withdraw struct {
	UserID      int64     `json:"user_id,omitempty" db:"user_id"`
	OrderID     string    `json:"order" db:"order_id"`
	Sum         float64   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}
