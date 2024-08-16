package models

import (
	"time"
)

type Order struct {
	Id               string
	User_id          string
	Product_id       string
	Quantity         int
	Size             string
	Status           string
	Purchase_id      string
	Invoice_url      string
	Delivery_option  string
	Delivery_address string
	Total_price      int
	Name             string
	Slug             string
	Description      string
	Price            int
	Image            string
	Stock            int
	Email            string

	Created_at *time.Time `json:"created_at,omitempty"`
	Updated_at *time.Time `json:"updated_at,omitempty"`
}

type Purchase struct {
	Id               []string
	User_id          string
	Product_id       string
	Quantity         int
	Size             string
	Status           string
	Purchase_id      string
	Invoice_url      string
	Delivery_option  string
	Delivery_Address string
	Total_price      int
}

type OrderItem struct {
	Name     string
	Quantity int
	Price    int
	ID       string
}

type Payment struct {
	ExternalID         string
	Amount             int
	PayerEmail         string
	Description        string
	SuccessRedirectURL string
	FailureRedirectURL string
}

type Dashboard struct {
	interval string
	data string
}
