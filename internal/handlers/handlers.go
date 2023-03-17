package handlers

import (
	"log"
	"net/http"
	"time"

	"gitlab.com/code-harbor/viwallet/internal/models"
	"gitlab.com/code-harbor/viwallet/internal/render"
	"gitlab.com/code-harbor/viwallet/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var repo repository.DatabaseRepository

func SetDBRepo(db repository.DatabaseRepository) {
	repo = db
}

// Home handler returns index.html page
func Home(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "index.page.tmpl")
	if err != nil {
		log.Println(err)
	}
}

// Login returns login page
func Login(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "login.page.tmpl")
	if err != nil {
		log.Println(err)
	}
}

// PostLogin gets the infomation from
func PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	role, username, err := repo.AuthenticateUser(email, password)
	if err != nil {
		log.Println("Login FAIL!")
		log.Println(err)
		return
	}
	log.Println("Login SUCCESS!")
	log.Println(role, username, err)
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "register.page.tmpl")
	if err != nil {
		log.Println(err)
	}
}

func PostRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/user/register", http.StatusTemporaryRedirect)
		return
	}

	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm_password")
	phone := r.Form.Get("phone")

	if confirmPassword != password {
		log.Println("password doesn't match")
		return
	}

	u := models.User{
		Email:       email,
		Username:    username,
		Password:    hashPassword(password),
		PhoneNumber: phone,
		Photo:       "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.CreateUser(u)

	if err != nil {
		log.Println(err)
	}
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashedPassword)
}
