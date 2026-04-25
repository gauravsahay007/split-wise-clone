package models

import "time"

type ExpenseShare struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}
type Expense struct {
	ID           int            `json:"id"`
	PaidBy       int            `json:"paid_by"` // user ID
	Amount       float64        `json:"amount"`
	GroupID      int            `json:"group_id"`
	UserIDs      []int          `json:"user_ids"` // participants
	Description  string         `json:"description"`
	Category     string         `json:"category"`
	ReceiptImage string         `json:"receipt_image"`
	SplitType    string         `json:"split_type"`
	CreatedAt    time.Time      `json:"created_at"`
	Shares       []ExpenseShare `json:"shares,omitempty"`
}
