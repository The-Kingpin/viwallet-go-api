package models

import "time"

type User struct {
	ID          int
	Email       string
	Password    string
	PhoneNumber string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Card struct {
	CardNumber     int
	ExpirationDate time.Time
	CVCNumber      int
	NamesOnCard    string
}
