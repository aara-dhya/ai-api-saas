package usage

import (
	"database/sql"
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
func (q *QuotaService) CheckQuota(apiKeyID string) (bool, error) {

	var used int

	err := q.db.QueryRow(
		`SELECT COALESCE(SUM(tokens_used), 0)
		 FROM usage_logs
		 WHERE api_key_id = $1
		 AND date_trunc('month', created_at) = date_trunc('month', NOW())`,
		apiKeyID,
	).Scan(&used)

	if err != nil {
		return false, err
	}

	if used >= q.limit {
		return false, nil
	}

	return true, nil
}
