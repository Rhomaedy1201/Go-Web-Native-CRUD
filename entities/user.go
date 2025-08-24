package entities

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
    ID       int    `json:"id" db:"id"`
    Name     string `json:"name" db:"name"`
    Email    string `json:"email" db:"email"`
    Password string `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Login struct {
    Email    string `json:"email" form:"email" validate:"required,email"`
    Password string `json:"password" form:"password" validate:"required,min=6"`
}

type Register struct {
    Name           string `json:"name" form:"name" validate:"required,min=2,max=50"`
    Email          string `json:"email" form:"email" validate:"required,email"`
    Password       string `json:"password" form:"password" validate:"required,min=6"`
    RepeatPassword string `json:"repeat_password" form:"repeat_password" validate:"required,eqfield=Password"`
}

// Global validator instance
var validate *validator.Validate

func init() {
    validate = validator.New()
}

// Method untuk validasi struct Login
func (l *Login) Validate() error {
    return validate.Struct(l)
}

// Method untuk validasi struct Register
func (r *Register) Validate() error {
    return validate.Struct(r)
}

// Method untuk convert Register ke User
func (r *Register) ToUser() User {
    return User{
        Name:     r.Name,
        Email:    r.Email,
        Password: r.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
    }
}