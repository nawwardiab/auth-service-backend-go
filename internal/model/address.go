package model

import "time"

type Address struct {
	ID        int       `json:"id"`
	UId       int       `json:"u_id"`
	Addr_1    string    `json:"addr_1"`
	Addr_2    string    `json:"addr_2"`
	Zip       string    `json:"zip"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}