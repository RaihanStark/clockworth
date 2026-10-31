package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RaihanStark/clockworth/internal/auth"
	"github.com/RaihanStark/clockworth/internal/config"
	"github.com/RaihanStark/clockworth/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	queries    *database.Queries
	oauth      *auth.GoogleOAuth
	cfg        *config.Config
	stateStore map[string]time.Time // In production, use Redis or similar
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

type UserResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Picture     string  `json:"picture"`
	HourlyWorth float64 `json:"hourly_worth"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type EmailLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) createSessionAndSetCookie(c echo.Context, userID uuid.UUID) (*database.Session, error) {
	sessionToken, err := auth.GenerateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	session, err := h.queries.CreateSession(c.Request().Context(), database.CreateSessionParams{
		UserID:    userID,
		Token:     sessionToken,
		ExpiresAt: auth.GetSessionExpiry(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
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

	return &session, nil
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || !strings.Contains(email, "@") {
		return echo.NewHTTPError(http.StatusBadRequest, "valid email is required")
	}

	if len(req.Password) < 6 {
		return echo.NewHTTPError(http.StatusBadRequest, "password must be at least 6 characters")
	}

	// Check if user already exists
	_, err := h.queries.GetUserByEmail(c.Request().Context(), email)
	if err == nil {
		return echo.NewHTTPError(http.StatusConflict, "user with this email already exists")
	} else if err != sql.ErrNoRows {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check existing user")
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	// Create user
	user, err := h.queries.CreateEmailUser(c.Request().Context(), database.CreateEmailUserParams{
		Email:        email,
		Name:         sql.NullString{String: strings.TrimSpace(req.Name), Valid: strings.TrimSpace(req.Name) != ""},
		PasswordHash: sql.NullString{String: passwordHash, Valid: true},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
	}

	// Create session and set cookie
	session, err := h.createSessionAndSetCookie(c, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token: session.Token,
		User: UserResponse{
			ID:          user.ID.String(),
			Email:       user.Email,
			Name:        user.Name.String,
			Picture:     user.Picture.String,
			HourlyWorth: user.HourlyWorth,
		},
	})
}

func (h *AuthHandler) EmailLogin(c echo.Context) error {
	var req EmailLoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email and password are required")
	}

	user, err := h.queries.GetUserByEmail(c.Request().Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to query user")
	}

	if !user.PasswordHash.Valid || user.PasswordHash.String == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "this account was registered with Google. Please sign in with Google.")
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash.String) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	}

	// Create session and set cookie
	session, err := h.createSessionAndSetCookie(c, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: session.Token,
		User: UserResponse{
			ID:          user.ID.String(),
			Email:       user.Email,
			Name:        user.Name.String,
			Picture:     user.Picture.String,
			HourlyWorth: user.HourlyWorth,
		},
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	if h.cfg.GoogleClientID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Google authentication is not configured on this server")
	}

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
	user, err := h.queries.GetUserByGoogleID(c.Request().Context(), sql.NullString{String: userInfo.ID, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new user
			user, err = h.queries.CreateUser(c.Request().Context(), database.CreateUserParams{
				Email:    userInfo.Email,
				Name:     sql.NullString{String: userInfo.Name, Valid: userInfo.Name != ""},
				Picture:  sql.NullString{String: userInfo.Picture, Valid: userInfo.Picture != ""},
				GoogleID: sql.NullString{String: userInfo.ID, Valid: true},
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
			}
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
		}
	}

	// Create session and set cookie
	session, err := h.createSessionAndSetCookie(c, user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create session")
	}

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
