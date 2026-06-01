package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *sql.DB
}

func getKeyPrefix(key string) string {
	if len(key) < 12 {
		return key
	}
	return key[:12]
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return "sk_live_" + hex.EncodeToString(bytes), nil
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (s *Service) CreateAPIKey(userID string, name string) (string, error) {

	// generate raw API key
	rawKey, err := generateAPIKey()
	if err != nil {
		return "", err
	}

	// hash API key
	hashedKey, err := bcrypt.GenerateFromPassword(
		[]byte(rawKey),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	// extract prefix
	prefix := getKeyPrefix(rawKey)

	query := `
	INSERT INTO api_keys (
		user_id,
		name,
		key,
		key_prefix,
		plan_id
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		(SELECT id FROM plans WHERE name = 'free')
	)
	`

	_, err = s.db.Exec(
		query,
		userID,
		name,
		string(hashedKey),
		prefix,
	)

	if err != nil {
		return "", err
	}

	// return RAW key only once
	return rawKey, nil
}
