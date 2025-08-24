package middlewares

import (
	"net/http"
	"strconv"
)

func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check login cookie
		loginCookie, err := r.Cookie("logged_in")
		if err != nil || loginCookie.Value != "true" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Check user_id cookie
		userIDCookie, err := r.Cookie("user_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate user_id
		if _, err := strconv.Atoi(userIDCookie.Value); err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// User authenticated, proceed
		next(w, r)
	}
}

func GuestOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check jika sudah login
		if loginCookie, err := r.Cookie("logged_in"); err == nil && loginCookie.Value == "true" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// User belum login, proceed
		next(w, r)
	}
}
