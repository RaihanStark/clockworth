package auth

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/RaihanStark/clockworth/internal/database"
	"github.com/labstack/echo/v4"
)

type contextKey string

const UserContextKey contextKey = "user"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func AuthMiddleware(queries *database.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// Check cookie as fallback
				cookie, err := c.Cookie("session_token")
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "missing authentication token")
				}
				authHeader = cookie.Value
			}

			// Remove "Bearer " prefix if present
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Get session from database
			session, err := queries.GetSessionByToken(c.Request().Context(), token)
			if err != nil {
				if err == sql.ErrNoRows {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
				}
				// Check if request was cancelled (e.g., user refreshed page)
				if err == context.Canceled {
					log.Printf("Auth middleware: request cancelled for token %s...", token[:min(8, len(token))])
					return echo.NewHTTPError(http.StatusRequestTimeout, "request cancelled")
				}
				// Log the actual error for debugging
				log.Printf("Auth middleware: database error for token %s...: %v", token[:min(8, len(token))], err)
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to validate session")
			}

			// Get user from database
			user, err := queries.GetUserByID(c.Request().Context(), session.UserID)
			if err != nil {
				if err == sql.ErrNoRows {
					return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
				}
				// Check if request was cancelled
				if err == context.Canceled {
					log.Printf("Auth middleware: request cancelled while fetching user")
					return echo.NewHTTPError(http.StatusRequestTimeout, "request cancelled")
				}
				// Log the actual error for debugging
				log.Printf("Auth middleware: failed to get user %s: %v", session.UserID, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
			}

			// Store user in context
			c.Set(string(UserContextKey), &user)

			return next(c)
		}
	}
}

func GetUserFromContext(c echo.Context) (*database.User, bool) {
	user, ok := c.Get(string(UserContextKey)).(*database.User)
	return user, ok
}
