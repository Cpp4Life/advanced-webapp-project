package model

type Slide struct {
	Id             uint     `json:"id,omitempty"`
	PresentationId uint     `json:"presentation_id,omitempty"`
	Type           uint     `json:"type,omitempty"`
	Content        *Content `json:"content,omitempty"`
}

type Content struct {
	Id        uint       `json:"id,omitempty"`
	SlideId   uint       `json:"slide_id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Meta      string     `json:"meta,omitempty"`
	Options   []*Option  `json:"options,omitempty"`
	Heading   *Heading   `json:"heading,omitempty"`
	Paragraph *Paragraph ` json:"paragraph,omitempty"`
}

type Option struct {
	Id         uint   `json:"id,omitempty"`
	Name       string `json:"name,omitempty" binding:"required"`
	Image      string `json:"image,omitempty"`
	ContentId  uint   `json:"content_id,omitempty"`
	TotalVotes uint   `json:"total_votes"`
}

type Heading struct {
	Id         uint   `json:"id,omitempty"`
	Heading    string `json:"heading,omitempty"`
	SubHeading string `json:"sub_heading,omitempty"`
	Image      string `json:"image,omitempty"`
	ContentId  uint   `json:"content_id,omitempty"`
}

type Paragraph struct {
	Id        uint   `json:"id,omitempty"`
	Heading   string `json:"heading,omitempty"`
	Text      string `json:"text,omitempty"`
	Image     string `json:"image,omitempty"`
	ContentId uint   `json:"content_id,omitempty"`
}
