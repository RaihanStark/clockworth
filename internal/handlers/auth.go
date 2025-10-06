package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/RaihanStark/clockworth/internal/auth"
	"github.com/RaihanStark/clockworth/internal/config"
	"github.com/RaihanStark/clockworth/internal/database"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	queries     *database.Queries
	oauth       *auth.GoogleOAuth
	cfg         *config.Config
	stateStore  map[string]time.Time // In production, use Redis or similar
}

func NewAuthHandler(queries *database.Queries, oauth *auth.GoogleOAuth, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		queries:    queries,
		oauth:      oauth,
		cfg:        cfg,
		stateStore: make(map[string]time.Time),
	}
}

type LoginResponse struct {
	URL string `json:"url"`
}

type CallbackResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	state, err := auth.GenerateStateToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate state token")
	}

	// Store state with expiration (5 minutes)
	h.stateStore[state] = time.Now().Add(5 * time.Minute)

	// Clean expired states
	go h.cleanExpiredStates()

	url := h.oauth.GetAuthURL(state)

	return c.JSON(http.StatusOK, LoginResponse{
		URL: url,
	})
}

func (h *AuthHandler) Callback(c echo.Context) error {
	state := c.QueryParam("state")
	code := c.QueryParam("code")

	// Validate state
	if expiry, ok := h.stateStore[state]; !ok || time.Now().After(expiry) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid or expired state")
	}
	delete(h.stateStore, state)

	// Exchange code for token
	token, err := h.oauth.Exchange(c.Request().Context(), code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to exchange code: %v", err))
	}

	// Get user info from Google
	userInfo, err := h.oauth.GetUserInfo(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to get user info: %v", err))
	}

	// Check if user exists
	user, err := h.queries.GetUserByGoogleID(c.Request().Context(), userInfo.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new user
			user, err = h.queries.CreateUser(c.Request().Context(), database.CreateUserParams{
				Email:    userInfo.Email,
				Name:     sql.NullString{String: userInfo.Name, Valid: userInfo.Name != ""},
				Picture:  sql.NullString{String: userInfo.Picture, Valid: userInfo.Picture != ""},
				GoogleID: userInfo.ID,
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
			}
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
		}
	}

	// Create session
	sessionToken, err := auth.GenerateSessionToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate session token")
	}

	session, err := h.queries.CreateSession(c.Request().Context(), database.CreateSessionParams{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: auth.GetSessionExpiry(),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

	// Set cookie (no SameSite for localhost development compatibility)
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", h.cfg.FrontendURL, session.Token)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	var token string

	if authHeader != "" {
		token = authHeader[len("Bearer "):]
	} else {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "no session found")
		}
		token = cookie.Value
	}

	// Delete session from database
	err := h.queries.DeleteSession(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to logout")
	}

	// Clear cookie
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":           user.ID,
		"email":        user.Email,
		"name":         user.Name.String,
		"picture":      user.Picture.String,
		"hourly_worth": user.HourlyWorth,
	})
}

func (h *AuthHandler) cleanExpiredStates() {
	now := time.Now()
	for state, expiry := range h.stateStore {
		if now.After(expiry) {
			delete(h.stateStore, state)
		}
	}
}
