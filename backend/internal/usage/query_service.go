package usage

import "database/sql"

type UsageSummary struct {
	TokensToday int          `json:"tokens_today"`
	CostToday   float64      `json:"cost_today"`
	TokensMonth int          `json:"tokens_month"`
	CostMonth   float64      `json:"cost_month"`
	Daily       []DailyUsage `json:"daily"`
}

type QueryService struct {
	db *sql.DB
}

func NewQueryService(db *sql.DB) *QueryService {
	return &QueryService{db: db}
}

func (q *QueryService) UsageSummary(apiKeyID string) (*UsageSummary, error) {

	summary := &UsageSummary{}

	// ---- TODAY ----
	err := q.db.QueryRow(
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
	err = q.db.QueryRow(
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

	// ---- LAST 7 DAYS ----
	rows, err := q.db.Query(
		`SELECT 
			DATE(created_at) as day,
			COALESCE(SUM(tokens_used), 0)
		FROM usage_logs
		WHERE api_key_id = $1
		AND created_at >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY day
		ORDER BY day`,
		apiKeyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d DailyUsage
		if err := rows.Scan(&d.Date, &d.Tokens); err != nil {
			return nil, err
		}
		summary.Daily = append(summary.Daily, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summary, nil
}
