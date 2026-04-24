package models

type Expense struct {
	ID      int     `json:"id"`
	PaidBy  int     `json:"paid_by"` // user ID
	Amount  float64 `json:"amount"`
	UserIDs []int   `json:"user_ids"` // participants
}
