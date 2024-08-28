package models

import (
	"time"
)

type TypeStat struct {
	Id               string
	User_id          string `json:"user_id"`
	Type             string `json:"type"`
	Easy_correct     int    `json:"easy_correct"`
	Easy_incorrect   int    `json:"easy_incorrect"`
	Medium_correct   int    `json:"medium_correct"`
	Medium_incorrect int    `json:"medium_incorrect"`
	Hard_correct     int    `json:"hard_correct"`
	Hard_incorrect   int    `json:"hard_incorrect"`

	Created_at *time.Time
	Updated_at *time.Time
}

type TotalStat struct {
	Total_score   int
	Highest_score int
	Rank          int
	Count         int
}
