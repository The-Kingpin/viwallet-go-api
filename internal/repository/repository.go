package repository

import (
	"gitlab.com/code-harbor/viwallet/internal/models"
)

type DatabaseRepository interface {
	GetUserByID(id int) (models.User, error)
	CreateUser(user models.User) error
	AuthenticateUser(email, password string) (int, string, error)
}
