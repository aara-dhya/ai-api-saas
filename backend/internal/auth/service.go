package auth

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// ----------------------
// SIGNUP
// ----------------------

func (s *Service) Signup(email, password string) (string, error) {

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	var userID string

	err = s.db.QueryRow(
		`INSERT INTO users (email, password)
		 VALUES ($1, $2)
		 RETURNING id`,
		email,
		string(hashed),
	).Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, nil
}

// ----------------------
// LOGIN
// ----------------------

func (s *Service) Login(email, password string) (string, error) {

	var userID string
	var hashedPassword string

	err := s.db.QueryRow(
		`SELECT id, password FROM users WHERE email = $1`,
		email,
	).Scan(&userID, &hashedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 🔥 generate JWT
	token, err := GenerateToken(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}
