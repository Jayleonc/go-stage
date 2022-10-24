package models

import (
	"time"
)

type User struct {
	ID           string `json:"id"`
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name" validate:"required,min=2,max=100"`
	Password     string `json:"password" validate:"required,min=6"`
	Email        string `json:"email" validate:"email,required"`
	Phone        string `json:"phone" validate:"required"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	UserType     string `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FilterUsr struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserPageDto struct {
	PageNum  int    `form:"page_num" json:"page_num"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"`
	Desc     bool   `form:"desc" json:"desc"`
	UserName string `form:"user_name" json:"user_name"`
	Email    string `form:"email" json:"email" validate:"email"`
	Phone    string `form:"phone" json:"phone" `
}
