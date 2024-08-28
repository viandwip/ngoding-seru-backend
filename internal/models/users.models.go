package models

import (
	"time"
)

type User struct {
	Id            string `form:"id"`
	Email         string `form:"email"`
	Password      string `form:"password" valid:"required, stringlength(6|100)~Password min. 6 chars"`
	Phone_number  string `form:"phone_number" valid:"required"`
	Role          string
	Image         string
	Address       string `form:"address"`
	Full_name     string `form:"full_name"`
	Birthday      string `form:"birthday"`
	Gender        string `form:"gender"`
	Highest_score string
	IsGoogle      bool

	Created_at *time.Time `form:"created_at,omitempty"`
	Updated_at *time.Time `form:"updated_at,omitempty"`
}
