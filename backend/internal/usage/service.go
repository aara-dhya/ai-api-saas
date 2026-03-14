package usage

import (
	"database/sql"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) LogUsage(apiKeyID string, model string, tokens int) error {

	_, err := s.db.Exec(`
		INSERT INTO usage_logs (api_key_id, model, tokens_used)
		VALUES ($1, $2, $3)
	`, apiKeyID, model, tokens)

	return err
}
