package handlers

import (
	"net/http"

	"github.com/RaihanStark/clockworth/internal/auth"
	"github.com/RaihanStark/clockworth/internal/database"
	"github.com/labstack/echo/v4"
)

type SettingsHandler struct {
	queries *database.Queries
}

func NewSettingsHandler(queries *database.Queries) *SettingsHandler {
	return &SettingsHandler{
		queries: queries,
	}
}

type UpdateHourlyWorthRequest struct {
	HourlyWorth float64 `json:"hourly_worth"`
}

type SettingsResponse struct {
	HourlyWorth float64 `json:"hourly_worth"`
}

func (h *SettingsHandler) GetSettings(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	return c.JSON(http.StatusOK, SettingsResponse{
		HourlyWorth: float64(user.HourlyWorth),
	})
}

func (h *SettingsHandler) UpdateHourlyWorth(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	var req UpdateHourlyWorthRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.HourlyWorth < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "hourly_worth must be a positive number")
	}

	updatedUser, err := h.queries.UpdateUserHourlyWorth(c.Request().Context(), database.UpdateUserHourlyWorthParams{
		ID:          user.ID,
		HourlyWorth: req.HourlyWorth,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update hourly worth")
	}

	return c.JSON(http.StatusOK, SettingsResponse{
		HourlyWorth: float64(updatedUser.HourlyWorth),
	})
}
