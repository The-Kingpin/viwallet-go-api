package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"gitlab.com/code-harbor/viwallet/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const TimeOut time.Duration = 3 * time.Second

func (pg *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User

	// query := ``
	return user, nil
}

func (pg *postgresDBRepo) CreateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	var newId int

	// Check if user is already registered with the given email
	err := pg.DB.QueryRowContext(ctx, "select id from users where email = $1", user.Email).Scan(&newId)

	if err == nil {
		return errors.New("user already exists")
	}

	query := `
		insert into users (email, username, password, phone_number, profile_photo, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)
	`
	err = pg.DB.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		user.Password,
		user.PhoneNumber,
		user.Photo,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&newId)

	if err != nil {
		return err
	} else {
		log.Println("user with id:", newId, "was created")
	}

	return nil
}

func (pg *postgresDBRepo) AuthenticateUser(email, inputPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()

	var id int
	var password string
	var username string

	// search for existing user with the given email. On success we get hashed password and id from the database
	err := pg.DB.QueryRowContext(ctx, "select id, password, username from users where email = $1", email).Scan(&id, &password, &username)

	if err != nil {
		return 0, "", errors.New("wrong credentials")
	}

	// compare the hashed password in the database with one provided in the login form
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(inputPassword))
	if err != nil {
		return id, "", err
	}

	return id, username, nil
}
