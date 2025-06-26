package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

var (
	ErrUserNotFound = errors.New("user not found in context")
	jwtSecret       = []byte("your-super-secret-jwt-key-change-in-production")
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}
type rateLimiter struct {
	visits map[string]int
	mu     sync.Mutex
	limit  int
	window time.Duration
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path, "status", rec.status, "duration", time.Since(start))
	})
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					"error", rec,
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visits: make(map[string]int),
		limit:  limit,
		window: window,
	}

	go func() {
		for {
			time.Sleep(rl.window)
			rl.mu.Lock()
			rl.visits = make(map[string]int)
			rl.mu.Unlock()
		}
	}()
	return rl
}

func (rl *rateLimiter) MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			slog.Warn("invalid RemoteAddr", "addr", r.RemoteAddr)
			http.Error(w, "Invalid address", http.StatusInternalServerError)
			return
		}
		ip := host

		rl.mu.Lock()
		rl.visits[ip]++
		count := rl.visits[ip]
		rl.mu.Unlock()

		if count >= rl.limit {
			slog.Warn("rate limit exceeded", "ip", ip, "count", count)
			http.Error(w, "Rate Limit Exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// optional: validate expiration, audience, etc.
		userID, ok := (*claims)["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext extracts the user ID from request context
func FromContext(ctx context.Context) (string, error) {
	userIDVal := ctx.Value(UserIDKey)
	if userIDVal == nil {
		return "", ErrUserNotFound
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return "", ErrUserNotFound
	}
	return userID, nil
}
