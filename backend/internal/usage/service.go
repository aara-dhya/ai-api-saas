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

func (s *Service) LogUsage(apiKeyID, model string, tokens int, cost float64) error {

	_, err := s.db.Exec(
		`INSERT INTO usage_logs (api_key_id, endpoint, tokens_used, cost)
		 VALUES ($1, $2, $3, $4)`,
		apiKeyID,
		model,
		tokens,
		cost,
	)

	return err
}
