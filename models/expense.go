package models

type Expense struct {
	ID      int     `json:"id"`
	PaidBy  int     `json:"paid_by"` // user ID
	Amount  float64 `json:"amount"`
	GroupID int     `json:"group_id"`
	UserIDs []int   `json:"user_ids"` // participants
}
