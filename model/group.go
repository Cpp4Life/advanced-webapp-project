package model

import (
	"time"
)

type Group struct {
	Id            uint           `json:"id,omitempty"`
	Name          string         `json:"name" binding:"required"`
	Link          string         `json:"link"`
	Desc          string         `json:"desc,omitempty"`
	CreatedAt     time.Time      `json:"created_at,omitempty"`
	Owner         *User          `json:"owner,omitempty"`
	GroupPresInfo *GroupPresInfo `json:"group_pres_info,omitempty"`
}

type GroupPresInfo struct {
	PresId  uint `json:"pres_id,omitempty"`
	GroupId uint `json:"group_id,omitempty"`
	UserId  uint `json:"user_id,omitempty"`
}
