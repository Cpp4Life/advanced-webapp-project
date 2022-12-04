package model

import (
	"time"
)

type Group struct {
	Id        uint      `json:"id,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Link      string    `json:"link" binding:"required"`
	Desc      string    `json:"desc,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Owner     *User     `json:"owner,omitempty"`
}
