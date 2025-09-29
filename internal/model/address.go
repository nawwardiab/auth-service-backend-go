package model

import "time"

type Address struct {
	ID int
	UId int
  Addr_1 string 
  Addr_2 string 
  Zip string 
  City string 
  Country string  
  IsDefault bool
  CreatedAt time.Time
  UpdatedAt time.Time
}