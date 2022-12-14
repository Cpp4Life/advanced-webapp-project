package model

import "time"

type GroupUser struct {
	GroupInfo *Group    `json:"group,omitempty"`
	UserInfo  *User     `json:"user_id,omitempty"`
	Role      string    `json:"role,omitempty"`
	JoinedAt  time.Time `json:"joined_at,omitempty"`
}

type Member struct {
	UserId uint   `json:"user_id,omitempty"`
	Email  string `json:"string"`
	Role   string `json:"role,omitempty"`
}
