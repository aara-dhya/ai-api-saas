package usage

import (
	"database/sql"
	"time"
)

type QuotaService struct {
	db *sql.DB
}

func NewQuotaService(db *sql.DB) *QuotaService {
	return &QuotaService{
		db: db,
	}
}

// CheckQuota verifies whether the API key has remaining quota for the current month.
func (q *QuotaService) CheckQuota(apiKeyID string) (bool, error) {

	// get start of current month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var usedTokens int

	err := q.db.QueryRow(
		`SELECT COALESCE(SUM(tokens),0)
		 FROM usage_logs
		 WHERE api_key_id = $1
		 AND created_at >= $2`,
		apiKeyID,
		startOfMonth,
	).Scan(&usedTokens)

	if err != nil {
		return false, err
	}

	// simple static quota (can later come from plans table)
	const monthlyQuota = 100000

	if usedTokens >= monthlyQuota {
		return false, nil
	}

	return true, nil
}
