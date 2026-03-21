package usage

import (
	"ai-api-saas/internal/billing"
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

// ----------------------
// USAGE SUMMARY
// ----------------------

type UsageSummary struct {
	TokensToday int     `json:"tokens_today"`
	CostToday   float64 `json:"cost_today"`
	TokensMonth int     `json:"tokens_month"`
	CostMonth   float64 `json:"cost_month"`
}

func (s *Service) UsageSummary(apiKeyID string) (*UsageSummary, error) {

	summary := &UsageSummary{}

	// ---- TODAY ----
	err := s.db.QueryRow(
		`SELECT 
			COALESCE(SUM(tokens_used), 0),
			COALESCE(SUM(cost), 0)
		FROM usage_logs
		WHERE api_key_id = $1
		AND created_at >= CURRENT_DATE
		AND created_at < CURRENT_DATE + INTERVAL '1 day'`,
		apiKeyID,
	).Scan(&summary.TokensToday, &summary.CostToday)

	if err != nil {
		return nil, err
	}

	// ---- THIS MONTH ----
	err = s.db.QueryRow(
		`SELECT 
			COALESCE(SUM(tokens_used), 0),
			COALESCE(SUM(cost), 0)
		FROM usage_logs
		WHERE api_key_id = $1
		AND created_at >= date_trunc('month', NOW())
		AND created_at < date_trunc('month', NOW()) + INTERVAL '1 month'`,
		apiKeyID,
	).Scan(&summary.TokensMonth, &summary.CostMonth)

	if err != nil {
		return nil, err
	}

	return summary, nil
}
