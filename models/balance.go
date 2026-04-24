package models

type Balance struct {
	FromUser int     `json:"from_user"`
	ToUser   int     `json:"to_user"`
	Amount   float64 `json:"amount"`
}
