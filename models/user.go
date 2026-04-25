package models

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Password   string `json:"-"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
}
