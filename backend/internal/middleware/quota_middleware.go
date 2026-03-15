package middleware

import (
	"fmt"
	"net/http"
)

// QuotaMiddleware temporarily bypasses quota for testing
type QuotaMiddleware struct {
	// quotaService *usage.QuotaService
	// We keep it here for later real implementation
}

func NewQuotaMiddleware( /*qs *usage.QuotaService*/ ) *QuotaMiddleware {
	return &QuotaMiddleware{
		// quotaService: qs,
	}
}

// Middleware bypasses quota check for testing
func (q *QuotaMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKeyID, ok := r.Context().Value(APIKeyIDKey).(string)
		if !ok {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		// TEMPORARY: bypass quota for testing
		fmt.Println("Quota check bypassed for API key:", apiKeyID)

		// proceed to next handler
		next.ServeHTTP(w, r)
	})
}
