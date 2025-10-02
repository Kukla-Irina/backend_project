package middleware

// каждому запросу присваивается уникальный request_id

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const RequestIDKey ctxKey = "request_id"
const HeaderRequestID = "X-Request-Id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(HeaderRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, rid)
		w.Header().Set(HeaderRequestID, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(RequestIDKey).(string)
	return v
}
