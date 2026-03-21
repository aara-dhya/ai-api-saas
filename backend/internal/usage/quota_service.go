package usage

import (
	"database/sql"
	"fmt"
)

const (
	DailyTokenLimit   = 10000
	MonthlyTokenLimit = 1000000
)

type QuotaService struct {
	db    *sql.DB
	limit int // monthly token limit
}

func NewQuotaService(db *sql.DB) *QuotaService {
	return &QuotaService{
		db:    db,
		limit: 1000000, // 🔥 set high for testing (1M tokens)
	}
}

// CheckQuota returns whether the API key is within its monthly quota
func (q *QuotaService) CheckQuota(apiKeyID string) error {

	var dailyUsed int
	var monthlyUsed int

	// ---- DAILY ----
	err := q.db.QueryRow(
		`SELECT COALESCE(SUM(tokens_used), 0)
		 FROM usage_logs
		 WHERE api_key_id = $1
		 AND created_at >= CURRENT_DATE
		 AND created_at < CURRENT_DATE + INTERVAL '1 day'`,
		apiKeyID,
	).Scan(&dailyUsed)

	if err != nil {
		return err
	}

	if dailyUsed >= DailyTokenLimit {
		return fmt.Errorf("daily quota exceeded")
	}

	// ---- MONTHLY ----
	err = q.db.QueryRow(
		`SELECT COALESCE(SUM(tokens_used), 0)
		 FROM usage_logs
		 WHERE api_key_id = $1
		 AND created_at >= date_trunc('month', NOW())
		 AND created_at < date_trunc('month', NOW()) + INTERVAL '1 month'`,
		apiKeyID,
	).Scan(&monthlyUsed)

	if err != nil {
		return err
	}

	if monthlyUsed >= MonthlyTokenLimit {
		return fmt.Errorf("monthly quota exceeded")
	}

	return nil
}
