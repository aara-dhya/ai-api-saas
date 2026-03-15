package middleware

import (
	"net/http"

	"ai-api-saas/internal/usage"
)

type QuotaMiddleware struct {
	quotaService *usage.QuotaService
}

func NewQuotaMiddleware(qs *usage.QuotaService) *QuotaMiddleware {
	return &QuotaMiddleware{
		quotaService: qs,
	}
}

func (q *QuotaMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKeyID, ok := r.Context().Value(APIKeyIDKey).(string)
		if !ok {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		allowed, err := q.quotaService.CheckQuota(apiKeyID)
		if err != nil {
			http.Error(w, "quota check failed", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "monthly quota exceeded", http.StatusPaymentRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}
