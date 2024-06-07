package model

import "time"

type Expenses struct {
	Id               string    `json:"id"`
	Date             time.Time `json:"date"`
	Amount           float64   `json:"amount"`
	Transaction_type string    `json:"transaction_type"`
	Balance          float64   `json:"balance"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"createdat"`
	UpdatedAt        time.Time `json:"updatedat"`
	User_Id          string    `json:"user_id"`
}
