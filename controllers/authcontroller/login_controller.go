package authcontroller

import (
	"crud-go/entities"
	"crud-go/models/usermodel"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type LoginData struct {
	Email        string
	Errors       map[string]string
	FlashMessage string
}

func Index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles(
		"views/layouts/base_auth.html",
		"views/auth/login.html",
	)
	if err != nil {
		panic(err)
	}

	data := LoginData{
		Errors: make(map[string]string),
	}

	// Check untuk flash message dari cookie
	if cookie, err := r.Cookie("flash_message"); err == nil {
		data.FlashMessage = cookie.Value
		// Hapus cookie setelah dibaca
		deleteCookie := &http.Cookie{
			Name:   "flash_message",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, deleteCookie)
	}

	temp.Execute(w, data)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form data
	loginData := entities.Login{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: r.FormValue("password"),
	}

	// Template data
	data := LoginData{
		Email:  loginData.Email,
		Errors: make(map[string]string),
	}

	// Validasi input
	if err := loginData.Validate(); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				data.Errors[field] = err.Field() + " tidak boleh kosong"
			case "email":
				data.Errors[field] = "Format email tidak valid"
			case "min":
				data.Errors[field] = err.Field() + " minimal " + err.Param() + " karakter"
			default:
				data.Errors[field] = "Input tidak valid"
			}
		}

		temp, _ := template.ParseFiles(
			"views/layouts/base_auth.html",
			"views/auth/login.html",
		)
		temp.Execute(w, data)
		return
	}

	// Check rate limiting
	if isRateLimited(r.RemoteAddr) {
		data.Errors["general"] = "Terlalu banyak percobaan login. Coba lagi dalam 5 menit."
		temp, _ := template.ParseFiles(
			"views/layouts/base_auth.html",
			"views/auth/login.html",
		)
		temp.Execute(w, data)
		return
	}

	// Get user dari database
	user, err := usermodel.GetByEmail(loginData.Email)
	if err != nil {
		// Jangan beri tahu apakah email tidak ditemukan (security)
		data.Errors["general"] = "Email atau password salah"
		recordFailedAttempt(r.RemoteAddr)
		
		temp, _ := template.ParseFiles(
			"views/layouts/base_auth.html",
			"views/auth/login.html",
		)
		temp.Execute(w, data)
		return
	}

	// Verify password
	if err := usermodel.VerifyPassword(user.Password, loginData.Password); err != nil {
		data.Errors["general"] = "Email atau password salah"
		recordFailedAttempt(r.RemoteAddr)
		
		temp, _ := template.ParseFiles(
			"views/layouts/base_auth.html",
			"views/auth/login.html",
		)
		temp.Execute(w, data)
		return
	}

	// Login berhasil - buat session
	if err := createUserSession(w, user); err != nil {
		data.Errors["general"] = "Gagal membuat session"
		temp, _ := template.ParseFiles(
			"views/layouts/base_auth.html",
			"views/auth/login.html",
		)
		temp.Execute(w, data)
		return
	}

	// Clear failed attempts
	clearFailedAttempts(r.RemoteAddr)

	// Redirect ke dashboard atau halaman utama
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Session management
func createUserSession(w http.ResponseWriter, user entities.User) error {
	// Set multiple cookies untuk security
	
	// User ID cookie (HttpOnly)
	userIDCookie := &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(user.ID),
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 7 hari
		HttpOnly: true,
		Secure:   false, // Set true jika menggunakan HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	// User name cookie (untuk display)
	userNameCookie := &http.Cookie{
		Name:     "user_name",
		Value:    user.Name,
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 7 hari
		HttpOnly: false, // Bisa diakses JavaScript untuk display
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	// Login status cookie
	loginCookie := &http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 7 hari
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, userIDCookie)
	http.SetCookie(w, userNameCookie)
	http.SetCookie(w, loginCookie)

	return nil
}

// Simple rate limiting (in production, use Redis atau database)
var failedAttempts = make(map[string][]time.Time)

func recordFailedAttempt(ip string) {
	now := time.Now()
	if attempts, exists := failedAttempts[ip]; exists {
		// Hapus attempt yang lebih dari 5 menit
		var validAttempts []time.Time
		for _, attempt := range attempts {
			if now.Sub(attempt) < 5*time.Minute {
				validAttempts = append(validAttempts, attempt)
			}
		}
		failedAttempts[ip] = append(validAttempts, now)
	} else {
		failedAttempts[ip] = []time.Time{now}
	}
}

func isRateLimited(ip string) bool {
	if attempts, exists := failedAttempts[ip]; exists {
		now := time.Now()
		var validAttempts []time.Time
		for _, attempt := range attempts {
			if now.Sub(attempt) < 5*time.Minute {
				validAttempts = append(validAttempts, attempt)
			}
		}
		failedAttempts[ip] = validAttempts
		return len(validAttempts) >= 5 // Max 5 attempts dalam 5 menit
	}
	return false
}

func clearFailedAttempts(ip string) {
	delete(failedAttempts, ip)
}

// Logout function
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Hapus semua cookies
	cookies := []string{"user_id", "user_name", "logged_in"}
	
	for _, cookieName := range cookies {
		cookie := &http.Cookie{
			Name:   cookieName,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}