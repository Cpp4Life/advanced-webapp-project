package model

import "time"

type GroupUser struct {
	GroupInfo *Group    `json:"group,omitempty"`
	UserInfo  *User     `json:"user_id,omitempty"`
	Role      string    `json:"role,omitempty"`
	JoinedAt  time.Time `json:"joined_at,omitempty"`
}
