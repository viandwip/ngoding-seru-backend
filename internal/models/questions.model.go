package models

import (
	"time"
)

type Question struct {
	Id          string
	Image       string
	Type        string `form:"type"`
	Level       string `form:"level"`
	Question    string `form:"question"`
	Option_a    string `form:"option_a"`
	Option_b    string `form:"option_b"`
	Option_c    string `form:"option_c"`
	Option_d    string `form:"option_d"`
	Explanation string `form:"explanation"`
	Answer      string `form:"answer"`

	Created_at *time.Time `form:"created_at,omitempty"`
	Updated_at *time.Time `form:"updated_at,omitempty"`
}
