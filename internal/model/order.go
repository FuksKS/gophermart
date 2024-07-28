package model

import "time"

type Order struct {
	Number     string    `json:"number" db:"order_id"`
	Status     string    `json:"status" db:"status"`
	Accrual    int64     `json:"getaccrual,omitempty" db:"getaccrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
