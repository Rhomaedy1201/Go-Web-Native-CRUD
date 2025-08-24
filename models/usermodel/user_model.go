package usermodel

import (
	"crud-go/config"
	"crud-go/entities"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Create user baru untuk registrasi
func Create(user entities.User) error {
    // Hash password sebelum disimpan
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %v", err)
    }

    // Query insert
    query := `
        INSERT INTO users (name, email, password, created_at, updated_at) 
        VALUES (?, ?, ?, ?, ?)
    `

    // Execute query
    _, err = config.DB.Exec(
        query,
        user.Name,
        user.Email,
        string(hashedPassword),
        time.Now(),
        time.Now(),
    )

    if err != nil {
        return fmt.Errorf("failed to create user: %v", err)
    }

    return nil
}

// Check apakah email sudah ada
func IsEmailExists(email string) (bool, error) {
    query := "SELECT COUNT(*) FROM users WHERE email = ?"
    
    var count int
    err := config.DB.QueryRow(query, email).Scan(&count)
    if err != nil {
        return false, fmt.Errorf("failed to check email: %v", err)
    }

    return count > 0, nil
}

// Get user by email untuk login
func GetByEmail(email string) (entities.User, error) {
    var user entities.User
    
    query := `
        SELECT id, name, email, password, created_at, updated_at 
        FROM users 
        WHERE email = ?
    `

    err := config.DB.QueryRow(query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return user, fmt.Errorf("user not found")
        }
        return user, fmt.Errorf("failed to get user: %v", err)
    }

    return user, nil
}

// Verify password untuk login
func VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}