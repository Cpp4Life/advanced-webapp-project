package model

type Slide struct {
	Id             uint     `json:"id,omitempty"`
	PresentationId uint     `json:"presentation_id,omitempty"`
	Type           uint     `json:"type,omitempty"`
	Content        *Content `json:"content,omitempty"`
}

type Content struct {
	Id      uint      `json:"id,omitempty"`
	SlideId uint      `json:"slide_id,omitempty"`
	Title   string    `json:"title,omitempty"`
	Meta    string    `json:"meta,omitempty"`
	Options []*Option `json:"options,omitempty"`
}

type Option struct {
	Id         uint   `json:"id,omitempty"`
	Name       string `json:"name,omitempty" binding:"required"`
	Image      string `json:"image,omitempty"`
	ContentId  uint   `json:"content_id,omitempty"`
	TotalVotes uint   `json:"total_votes,omitempty"`
}
