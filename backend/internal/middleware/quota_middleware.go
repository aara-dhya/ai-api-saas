package middleware

import (
	"net/http"
)

type QuotaMiddleware struct {
	checkQuota func(string) (bool, error)
}

func NewQuotaMiddleware(check func(string) (bool, error)) *QuotaMiddleware {
	return &QuotaMiddleware{
		checkQuota: check,
	}
}

func (q *QuotaMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKeyID, ok := r.Context().Value(APIKeyIDKey).(string)
		if !ok {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		allowed, err := q.checkQuota(apiKeyID)
		if err != nil {
			http.Error(w, "quota check failed", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "quota exceeded", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
