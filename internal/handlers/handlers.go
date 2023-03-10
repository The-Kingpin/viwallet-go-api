package handlers

import (
	"log"
	"net/http"

	"gitlab.com/code-harbor/viwallet/internal/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "index.page.tmpl")
	if err != nil {
		log.Fatal(err)
	}
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "login.page.tmpl")

	if err != nil {
		log.Fatal(err)
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request) {

}

func GetRegister(w http.ResponseWriter, r *http.Request) {
	err := render.RenderTemplate(w, r, "register.page.tmpl")
	if err != nil {
		log.Fatal(err)
	}
}

func PostRegister(w http.ResponseWriter, r *http.Request) {

}
