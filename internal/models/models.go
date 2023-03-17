package models

import (
	"time"
)

type User struct {
	ID          int
	Email       string
	Username    string
	Password    string
	PhoneNumber string
	Photo       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Card struct {
	ID             int
	CardNumber     string
	ExpirationDate time.Time
	CVCNumber      int
	NamesOnCard    string
}
