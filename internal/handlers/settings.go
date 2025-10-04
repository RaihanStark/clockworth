package handlers

import (
	"database/sql"
	"net/http"

	"github.com/RaihanStark/clockworth/internal/auth"
	"github.com/RaihanStark/clockworth/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
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
	HourlyWorth float64 `json:"hourlyWorth" validate:"required,min=0"`
}

type HourlyWorthResponse struct {
	HourlyWorth float64 `json:"hourlyWorth"`
}

func (h *SettingsHandler) GetHourlyWorth(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	// Get hourly worth from database
	hourlyWorth, err := h.queries.GetUserHourlyWorth(c.Request().Context(), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default value if not set
			return c.JSON(http.StatusOK, HourlyWorthResponse{
				HourlyWorth: 25.00,
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get hourly worth")
	}

	// Convert decimal to float64
	dec, err := decimal.NewFromString(hourlyWorth)
	if err != nil {
		// Return default value if parsing fails
		return c.JSON(http.StatusOK, HourlyWorthResponse{
			HourlyWorth: 25.00,
		})
	}
	value, _ := dec.Float64()

	return c.JSON(http.StatusOK, HourlyWorthResponse{
		HourlyWorth: value,
	})
}

func (h *SettingsHandler) UpdateHourlyWorth(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	// Bind request
	var req UpdateHourlyWorthRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update hourly worth in database
	decimalValue := decimal.NewFromFloat(req.HourlyWorth).String()
	_, err := h.queries.UpdateUserHourlyWorth(c.Request().Context(), database.UpdateUserHourlyWorthParams{
		ID:          user.ID,
		HourlyWorth: decimalValue,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update hourly worth")
	}

	return c.JSON(http.StatusOK, HourlyWorthResponse{
		HourlyWorth: req.HourlyWorth,
	})
}