package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.com/code-harbor/viwallet/internal/config"
	"gitlab.com/code-harbor/viwallet/internal/driver"
	"gitlab.com/code-harbor/viwallet/internal/models"
	"gitlab.com/code-harbor/viwallet/internal/render"
	"gitlab.com/code-harbor/viwallet/internal/repository"
	"gitlab.com/code-harbor/viwallet/internal/repository/dbrepo"
	"golang.org/x/crypto/bcrypt"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepository
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func SetRepo(r *Repository) {
	Repo = r
}

// Home handler returns index.html page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("r.RemoteAddr: %v\n", r.ParseForm())

	username, ok := m.App.Session.Get(r.Context(), "username").(string)
	if !ok {
		username = ""
	}
	m.App.InfoLog.Println(username)
	err := render.Template(w, r, "index.page.tmpl", &models.TemplateData{
		Username: username,
	})
	if err != nil {
		log.Println(err)
	}
}

// Login returns login page
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "login.page.tmpl", &models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// PostLogin gets the infomation from
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	_, username, err := Repo.DB.AuthenticateUser(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

		return
	}

	m.App.InfoLog.Println("***", username)
	m.App.Session.Put(r.Context(), "username", username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "register.page.tmpl", &models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
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

	err = Repo.DB.CreateUser(u)

	if err != nil {
		log.Println(err)
	}
}

func (m *Repository) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// get user id from session
	// currentUserData := repo.GetUserByID(app.session.get(user.id))

	email := r.Form.Get("email")
	phoneNumber := r.Form.Get("phone")
	photo := r.Form.Get("photo")

	// email validation before assign the value
	// phone validation before assign the value
	// validation for format, size etc
	u := models.User{
		Email:       email,
		PhoneNumber: phoneNumber,
		Photo:       photo,
	}

	log.Println(u)

}

func (m *Repository) Trasnsactions(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "transactions.page.tmpl", &models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

func (m *Repository) Cards(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "cards.page.tmpl", &models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

func (m *Repository) Wallets(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("r.RemoteAddr: %v\n", r.RemoteAddr)
	err := render.Template(w, r, "wallets.page.tmpl", &models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hashedPassword)
}
