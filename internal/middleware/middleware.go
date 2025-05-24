package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/LikhithMar14/workout-tracker/internal/auth"
	"github.com/LikhithMar14/workout-tracker/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is a type for context keys to avoid conflicts
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userID"
	// UserEmailKey is the context key for user email
	UserEmailKey ContextKey = "userEmail"
)

// Middleware is a struct that holds dependencies for middleware functions
type Middleware struct {
	Logger        *log.Logger
	Authenticator auth.Authenticator
}

// NewMiddleware creates a new middleware instance
func NewMiddleware(logger *log.Logger, authenticator auth.Authenticator) *Middleware {
	return &Middleware{
		Logger:        logger,
		Authenticator: authenticator,
	}
}

// RequireAuth is a middleware that validates JWT tokens
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "missing authorization header"})
			return
		}

		// Check if the header starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header format"})
			return
		}

		// Extract the token
		tokenString := parts[1]
		if tokenString == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "missing token"})
			return
		}

		// Validate the token
		token, err := m.Authenticator.ValidateToken(tokenString)
		if err != nil {
			m.Logger.Printf("ERROR: token validation failed: %v", err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or expired token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token claims"})
			return
		}

		// Extract user information from claims
		userID, ok := claims["user_id"].(float64) // JWT numeric claims are float64
		if !ok {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid user ID in token"})
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid email in token"})
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), UserIDKey, int(userID))
		ctx = context.WithValue(ctx, UserEmailKey, email)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORS middleware to handle cross-origin requests
func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // In production, specify allowed origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequestLogger middleware to log requests
func (m *Middleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Log the request
		duration := time.Since(start)
		m.Logger.Printf("%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// ContentType middleware to set JSON content type for API responses
func (m *Middleware) ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// RecoverPanic middleware to recover from panics and return a 500 error
func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.Printf("PANIC: %v", err)
				utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetUserEmailFromContext extracts user email from request context
func GetUserEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(UserEmailKey).(string)
	if !ok {
		return "", fmt.Errorf("user email not found in context")
	}
	return email, nil
}

// RateLimiter is a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// RateLimit middleware to limit requests per IP
func (m *Middleware) RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			now := time.Now()

			// Clean old requests
			if times, exists := limiter.requests[ip]; exists {
				var validTimes []time.Time
				for _, t := range times {
					if now.Sub(t) < limiter.window {
						validTimes = append(validTimes, t)
					}
				}
				limiter.requests[ip] = validTimes
			}

			// Check rate limit
			if len(limiter.requests[ip]) >= limiter.limit {
				utils.WriteJSON(w, http.StatusTooManyRequests, utils.Envelope{"error": "rate limit exceeded"})
				return
			}

			// Add current request
			limiter.requests[ip] = append(limiter.requests[ip], now)

			next.ServeHTTP(w, r)
		})
	}
}
