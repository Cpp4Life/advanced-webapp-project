package model

import "time"

type User struct {
	Id            uint      `json:"id"`
	FullName      string    `json:"full_name,omitempty"`
	Password      string    `json:"password,omitempty"`
	Email         string    `json:"email,omitempty"`
	SavedPassword string    `json:"-"`
	Username      string    `json:"username,omitempty"`
	Address       string    `json:"address,omitempty"`
	ProfileImg    string    `json:"profile_img,omitempty"`
	UserTel       string    `json:"user_tel,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
