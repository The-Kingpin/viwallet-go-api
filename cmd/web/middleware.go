package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// CSRFProtect adds CSRF protection to all POST requests
func CSRFProtect(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		// make configuration for dynamicly setting the Secure field to true or false
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})	

	return csrfHandler
}

// Auth checks the requester is authenticated
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !Authenticated(r) {
			session.Put(r.Context(), "error", "Login first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Authenticated returns true, if user is authenticated and if not false
func Authenticated(r *http.Request) bool {
	exists := session.Exists(r.Context(), "user_id")
	return exists
}
