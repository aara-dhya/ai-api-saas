package usage

import (
	"ai-api-saas/internal/billing"
	"database/sql"
)

type DailyUsage struct {
	Date   string `json:"date"`
	Tokens int    `json:"tokens"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

// ----------------------
// USAGE LOGGING
// ----------------------

func (s *Service) LogUsage(apiKeyID, model string, tokens int) error {

	cost := billing.CalculateCost(model, tokens)

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
