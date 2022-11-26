package model

import "time"

type User struct {
	Id            uint      `json:"id"`
	FullName      string    `json:"full_name,omitempty" binding:"required"`
	Password      string    `json:"password,omitempty" binding:"required,gte=8"`
	Email         string    `json:"email,omitempty" binding:"required,email"`
	SavedPassword string    `json:"-"`
	Username      string    `json:"username,omitempty"`
	Address       string    `json:"address,omitempty"`
	ProfileImg    []byte    `json:"profile_img,omitempty"`
	UserTel       string    `json:"user_tel,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
