package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func generateKey() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return "sk_" + hex.EncodeToString(bytes), nil
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (s *Service) CreateAPIKey(userID string, name string) (string, error) {

	key, err := generateKey()
	if err != nil {
		return "", err
	}

	hash := hashKey(key)

	query := `
	INSERT INTO api_keys (user_id, key_hash, name)
	VALUES ($1, $2, $3)
	`

	_, err = s.db.Exec(query, userID, hash, name)
	if err != nil {
		return "", err
	}

	return key, nil
}
