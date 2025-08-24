package authcontroller

import (
	"crud-go/entities"
	"crud-go/models/usermodel"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RegisterData struct {
    Name   string
    Email  string
    Errors map[string]string
}

func IndexRegis(w http.ResponseWriter, r *http.Request){
	temp, err := template.ParseFiles(
		"views/layouts/base_auth.html",
		"views/auth/register.html",
	)
	if err != nil {
		panic(err)
	}

	// Kirim data kosong dengan struktur yang sama
	data := RegisterData{
		Errors: make(map[string]string),
	}
	temp.Execute(w, data)
}

func RegisterUser(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	registerData := entities.Register{
		Name: 		 r.FormValue("name"),
		Email: 		 r.FormValue("email"),
		Password: 	 r.FormValue("password"),
		RepeatPassword: r.FormValue("repeat_password"),
	}

	data := RegisterData{
		Name:   registerData.Name,
		Email:  registerData.Email,
		Errors: make(map[string]string),
	}

	if err := registerData.Validate(); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
            field := strings.ToLower(err.Field())
            switch err.Tag() {
            case "required":
                data.Errors[field] = err.Field() + " tidak boleh kosong"
            case "email":
                data.Errors[field] = "Format email tidak valid"
            case "min":
                data.Errors[field] = err.Field() + " minimal " + err.Param() + " karakter"
            case "max":
                data.Errors[field] = err.Field() + " maksimal " + err.Param() + " karakter"
            case "eqfield":
                data.Errors[field] = "Password tidak sama"
            default:
                data.Errors[field] = "Input tidak valid"
            }
        }

		temp, err := template.ParseFiles(
            "views/layouts/base_auth.html",
            "views/auth/register.html",
        )
        if err != nil {
            panic(err)
        }
        temp.Execute(w, data)
        return
	}

	// Check apakah email sudah ada
    emailExists, err := usermodel.IsEmailExists(registerData.Email)
    if err != nil {
        data.Errors["email"] = "Terjadi kesalahan saat memeriksa email"
        temp, _ := template.ParseFiles(
            "views/layouts/base_auth.html",
            "views/auth/register.html",
        )
        temp.Execute(w, data)
        return
    }

    if emailExists {
        data.Errors["email"] = "Email sudah terdaftar"
        temp, _ := template.ParseFiles(
            "views/layouts/base_auth.html",
            "views/auth/register.html",
        )
        temp.Execute(w, data)
        return
    }

    // Convert ke User dan simpan ke database
    user := registerData.ToUser()
    if err := usermodel.Create(user); err != nil {
        data.Errors["general"] = "Gagal mendaftarkan user"
        temp, _ := template.ParseFiles(
            "views/layouts/base_auth.html",
            "views/auth/register.html",
        )
        temp.Execute(w, data)
        return
    }

    // Set cookie untuk success message
    cookie := &http.Cookie{
        Name:     "flash_message",
        Value:    "Registrasi berhasil! Silakan login dengan akun Anda.",
        Path:     "/",
        MaxAge:   10, // 10 detik
        HttpOnly: true,
    }
    http.SetCookie(w, cookie)

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}