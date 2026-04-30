package models

type Friend struct {
	ID           int `json:"id"`
	UserID       int `json:"user_id"`
	FriendUserID int `json:"friend_user_id"`
}
