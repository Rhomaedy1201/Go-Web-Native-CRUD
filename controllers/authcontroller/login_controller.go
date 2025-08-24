package authcontroller

import (
	"html/template"
	"net/http"
)

type LoginData struct {
    Email        string
    Errors       map[string]string
    FlashMessage string
}

func Index(w http.ResponseWriter, r *http.Request){
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