package model

type User struct {
	Id         uint   `json:"id"`
	FullName   string `json:"full_name" validate:"required"`
	Password   string `json:"password" validate:"required,gte=8"`
	Username   string `json:"username"`
	Email      string `json:"email" validate:"email"`
	Address    string `json:"address,omitempty"`
	ProfileImg []byte `json:"profile_img,omitempty"`
	UserTel    string `json:"user_tel,omitempty"`
}
