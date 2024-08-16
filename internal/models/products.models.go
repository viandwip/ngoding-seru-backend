package models

import (
	"time"
)

type Product struct {
	Id          string
	Image       string
	Name        string `form:"name"`
	Description string `form:"description"`
	Price       int    `form:"price"`
	Stock       int    `form:"stock"`
	Size        string `form:"size"`
	Slug        string

	Created_at *time.Time `form:"created_at,omitempty"`
	Updated_at *time.Time `form:"updated_at,omitempty"`
}
