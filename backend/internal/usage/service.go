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

// ----------------------
// USAGE LOGGING
// ----------------------

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

// ----------------------
// COST CALCULATION
// ----------------------

func calculateCost(model string, tokens int) float64 {
	switch model {
	case "llama-3.1-8b-instant":
		return float64(tokens) * 0.0000002
	default:
		return 0
	}
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
		AND DATE(created_at) = CURRENT_DATE`,
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
		AND date_trunc('month', created_at) = date_trunc('month', NOW())`,
		apiKeyID,
	).Scan(&summary.TokensMonth, &summary.CostMonth)

	if err != nil {
		return nil, err
	}

	return summary, nil
}
