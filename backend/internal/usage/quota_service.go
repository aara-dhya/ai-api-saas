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

	var dailyLimit int
	var monthlyLimit int

	// ---- FETCH PLAN LIMITS ----
	err := q.db.QueryRow(
		`SELECT p.daily_token_limit, p.monthly_token_limit
		 FROM api_keys ak
		 JOIN plans p ON ak.plan_id = p.id
		 WHERE ak.id = $1`,
		apiKeyID,
	).Scan(&dailyLimit, &monthlyLimit)

	if err != nil {
		return err
	}

	// ---- DAILY USAGE ----
	err = q.db.QueryRow(
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

	if dailyUsed >= dailyLimit {
		return fmt.Errorf("daily quota exceeded")
	}

	// ---- MONTHLY USAGE ----
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

	if monthlyUsed >= monthlyLimit {
		return fmt.Errorf("monthly quota exceeded")
	}

	return nil
}
