package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RaihanStark/clockworth/internal/auth"
	"github.com/RaihanStark/clockworth/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ActivitiesHandler struct {
	queries *database.Queries
}

func NewActivitiesHandler(queries *database.Queries) *ActivitiesHandler {
	return &ActivitiesHandler{
		queries: queries,
	}
}

type CreateActivityRequest struct {
	Title        string `json:"title"`
	StartTime    string `json:"start_time"` // HH:MM format
	EndTime      string `json:"end_time"`   // HH:MM format
	Category     string `json:"category"`
	ROI          string `json:"roi"`
	ActivityDate string `json:"activity_date"` // YYYY-MM-DD format
}

type UpdateActivityRequest struct {
	Title     string `json:"title"`
	StartTime string `json:"start_time"` // HH:MM format
	EndTime   string `json:"end_time"`   // HH:MM format
	Category  string `json:"category"`
	ROI       string `json:"roi"`
}

type ActivityResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Category     string `json:"category"`
	ROI          string `json:"roi"`
	ActivityDate string `json:"activity_date"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func activityToResponse(activity database.Activity) ActivityResponse {
	return ActivityResponse{
		ID:           activity.ID.String(),
		Title:        activity.Title,
		StartTime:    activity.StartTime.Format("15:04"),
		EndTime:      activity.EndTime.Format("15:04"),
		Category:     activity.Category,
		ROI:          activity.Roi,
		ActivityDate: activity.ActivityDate.Format("2006-01-02"),
		CreatedAt:    activity.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    activity.UpdatedAt.Time.Format(time.RFC3339),
	}
}

// ListActivities handles GET /api/activities?date=YYYY-MM-DD
func (h *ActivitiesHandler) ListActivities(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	dateStr := c.QueryParam("date")
	if dateStr == "" {
		// Default to today
		dateStr = time.Now().Format("2006-01-02")
	}

	activityDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
	}

	activities, err := h.queries.GetActivitiesByUserAndDate(c.Request().Context(), database.GetActivitiesByUserAndDateParams{
		UserID:       user.ID,
		ActivityDate: activityDate,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch activities")
	}

	response := make([]ActivityResponse, len(activities))
	for i, activity := range activities {
		response[i] = activityToResponse(activity)
	}

	return c.JSON(http.StatusOK, response)
}

// GetActivity handles GET /api/activities/:id
func (h *ActivitiesHandler) GetActivity(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity ID")
	}

	activity, err := h.queries.GetActivityByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "activity not found")
	}

	// Verify the activity belongs to the user
	if activity.UserID != user.ID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	return c.JSON(http.StatusOK, activityToResponse(activity))
}

// CreateActivity handles POST /api/activities
func (h *ActivitiesHandler) CreateActivity(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	var req CreateActivityRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate required fields
	if req.Title == "" || req.StartTime == "" || req.EndTime == "" || req.Category == "" || req.ROI == "" || req.ActivityDate == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "all fields are required")
	}

	// Parse times
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start_time format, expected HH:MM")
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid end_time format, expected HH:MM")
	}

	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity_date format, expected YYYY-MM-DD")
	}

	fmt.Println("startTime", startTime)
	fmt.Println("endTime", endTime)
	fmt.Println("activityDate", activityDate)

	// Check for time overlap
	hasOverlap, err := h.queries.CheckActivityOverlap(c.Request().Context(), database.CheckActivityOverlapParams{
		UserID:       user.ID,
		ActivityDate: activityDate,
		Column3:      endTime,   // $3 = new_end_time for overlap check
		Column4:      startTime, // $4 = new_start_time for overlap check
	})
	if err != nil {
		fmt.Println("err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check activity overlap")
	}

	if hasOverlap {
		return echo.NewHTTPError(http.StatusConflict, "activity time overlaps with an existing activity")
	}

	fmt.Println("hasOverlap", hasOverlap)

	activity, err := h.queries.CreateActivity(c.Request().Context(), database.CreateActivityParams{
		UserID:       user.ID,
		Title:        req.Title,
		StartTime:    startTime,
		EndTime:      endTime,
		Category:     req.Category,
		Roi:          req.ROI,
		ActivityDate: activityDate,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create activity")
	}

	return c.JSON(http.StatusCreated, activityToResponse(activity))
}

// UpdateActivity handles PUT /api/activities/:id
func (h *ActivitiesHandler) UpdateActivity(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity ID")
	}

	// Verify the activity exists and belongs to the user
	existingActivity, err := h.queries.GetActivityByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "activity not found")
	}

	if existingActivity.UserID != user.ID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	var req UpdateActivityRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// Validate required fields
	if req.Title == "" || req.StartTime == "" || req.EndTime == "" || req.Category == "" || req.ROI == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "all fields are required")
	}

	// Parse times
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start_time format, expected HH:MM")
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid end_time format, expected HH:MM")
	}

	// Check for time overlap with other activities (excluding this one)
	hasOverlap, err := h.queries.CheckActivityOverlapExcludingID(c.Request().Context(), database.CheckActivityOverlapExcludingIDParams{
		UserID:       user.ID,
		ActivityDate: existingActivity.ActivityDate,
		ID:           id,
		Column4:      endTime,   // $4 = new_end_time for overlap check
		Column5:      startTime, // $5 = new_start_time for overlap check
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to check activity overlap")
	}

	if hasOverlap {
		return echo.NewHTTPError(http.StatusConflict, "activity time overlaps with an existing activity")
	}

	activity, err := h.queries.UpdateActivity(c.Request().Context(), database.UpdateActivityParams{
		ID:        id,
		Title:     req.Title,
		StartTime: startTime,
		EndTime:   endTime,
		Category:  req.Category,
		Roi:       req.ROI,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update activity")
	}

	return c.JSON(http.StatusOK, activityToResponse(activity))
}

// DeleteActivity handles DELETE /api/activities/:id
func (h *ActivitiesHandler) DeleteActivity(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid activity ID")
	}

	// Verify the activity exists and belongs to the user
	existingActivity, err := h.queries.GetActivityByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "activity not found")
	}

	if existingActivity.UserID != user.ID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	err = h.queries.DeleteActivity(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete activity")
	}

	return c.NoContent(http.StatusNoContent)
}
