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

// LogUsage logs API usage and calculates cost internally
func (s *Service) LogUsage(apiKeyID, model string, tokens int) error {

	cost := calculateCost(model, tokens)

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

// calculateCost determines cost based on model + tokens
func calculateCost(model string, tokens int) float64 {
	switch model {
	case "llama-3.1-8b-instant":
		return float64(tokens) * 0.0000002 // adjust pricing later
	default:
		return 0
	}
}

// UsageThisMonth returns total tokens used this month
func (s *Service) UsageThisMonth(apiKeyID string) (int, error) {

	var total int

	err := s.db.QueryRow(
		`SELECT COALESCE(SUM(tokens_used), 0)
		 FROM usage_logs
		 WHERE api_key_id = $1
		 AND date_trunc('month', created_at) = date_trunc('month', NOW())`,
		apiKeyID,
	).Scan(&total)

	return total, err
}
