package service

import (
    "database/sql"
    /"errors"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
    return &AuthService{db: db}
}

func (s *AuthService) Register(username, password string) error {
    // Хеширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // Сохранение пользователя в базе данных
    _, err = s.db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, hashedPassword)
    if err != nil {
        return err
    }

    return nil
}

func (s *AuthService) Login(username, password string) (bool, error) {
    var storedHash string
    err := s.db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&storedHash)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil // Пользователь не найден
        }
        return false, err
    }

    // Проверка пароля
    if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
        return false, nil // Неверный пароль
    }

    return true, nil
}