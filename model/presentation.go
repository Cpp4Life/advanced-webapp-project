package model

import "time"

type Pres struct {
	Id         uint      `json:"id,omitempty"`
	Name       string    `json:"name,omitempty" binding:"required"`
	Owner      *User     `json:"owner,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
