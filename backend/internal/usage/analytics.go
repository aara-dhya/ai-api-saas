package usage

func (s *Service) MonthlyUsage(apiKeyID string) (int, error) {

	var tokens int

	err := s.db.QueryRow(
		`SELECT COALESCE(SUM(tokens),0)
		 FROM usage_logs
		 WHERE api_key_id=$1
		 AND date_trunc('month', created_at)=date_trunc('month', now())`,
		apiKeyID,
	).Scan(&tokens)

	return tokens, err
}
