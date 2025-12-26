package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"event_planner_backend/config"
	"event_planner_backend/middleware"
	"event_planner_backend/models"
	"event_planner_backend/utils"
)

// CreateEventRequest represents the payload for creating an event.
type CreateEventRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description"`
	Location    string `json:"location" binding:"required,min=1,max=255"`
	EventDate   string `json:"eventDate" binding:"required"` // Format: "2006-01-02"
	EventTime   string `json:"eventTime" binding:"required"` // Format: "15:04:05" or "15:04"
}

// InviteUserRequest represents the payload for inviting a user to an event.
type InviteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"omitempty,oneof=organizer attendee"` // Default: attendee
}

// CreateEvent creates a new event and marks the creator as organizer.
func CreateEvent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	// Parse event date
	eventDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid eventDate format. Use YYYY-MM-DD: "+err.Error())
		return
	}

	// Parse event time (support both HH:MM:SS and HH:MM formats)
	eventTime := req.EventTime
	if len(eventTime) == 5 && eventTime[2] == ':' {
		eventTime = eventTime + ":00" // Add seconds if missing
	}
	_, err = time.Parse("15:04:05", eventTime)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid eventTime format. Use HH:MM:SS or HH:MM: "+err.Error())
		return
	}

	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		EventDate:   eventDate,
		EventTime:   eventTime,
		CreatedBy:   userID,
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create event
	if err := tx.Create(event).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "failed to create event: "+err.Error())
		return
	}

	// Mark creator as organizer
	attendee := &models.EventAttendee{
		EventID:   event.ID,
		UserID:    userID,
		Role:      "organizer",
		Status:    "going", // Creator is automatically going
		InvitedAt: time.Now(),
	}

	if err := tx.Create(attendee).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "failed to create attendee record: "+err.Error())
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"id":          event.ID,
		"title":       event.Title,
		"description": event.Description,
		"location":    event.Location,
		"eventDate":   event.EventDate.Format("2006-01-02"),
		"eventTime":   event.EventTime,
		"createdBy":   event.CreatedBy,
		"createdAt":   event.CreatedAt,
	})
}

// GetMyOrganizedEvents returns all events organized by the authenticated user.
func GetMyOrganizedEvents(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	var events []models.Event
	if err := config.DB.Where("created_by = ?", userID).
		Preload("Attendees", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch events: "+err.Error())
		return
	}

	// Format response
	response := make([]gin.H, len(events))
	for i, event := range events {
		response[i] = formatEventResponse(event)
	}

	c.JSON(http.StatusOK, response)
}

// GetMyInvitedEvents returns all events the authenticated user is invited to (as attendee or organizer).
func GetMyInvitedEvents(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	var attendees []models.EventAttendee
	if err := config.DB.Where("user_id = ?", userID).
		Preload("Event", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Organizer").Preload("Attendees", func(db *gorm.DB) *gorm.DB {
				return db.Preload("User")
			})
		}).
		Preload("User").
		Order("invited_at DESC").
		Find(&attendees).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch events: "+err.Error())
		return
	}

	// Format response with role and status
	response := make([]gin.H, len(attendees))
	for i, attendee := range attendees {
		event := attendee.Event
		response[i] = formatEventResponse(event)
		response[i]["role"] = attendee.Role
		response[i]["status"] = attendee.Status
		response[i]["invitedAt"] = attendee.InvitedAt
	}

	c.JSON(http.StatusOK, response)
}

// InviteUserToEvent invites a user to an event.
func InviteUserToEvent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		utils.JSONError(c, http.StatusBadRequest, "event ID required")
		return
	}

	var req InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	// Set default role
	role := req.Role
	if role == "" {
		role = "attendee"
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	// Use transaction for atomicity
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify event exists and user is organizer in one query
	var event models.Event
	if err := tx.Where("event_id = ? AND created_by = ?", eventID, userID).
		First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Check if user is organizer via attendees table
			var attendee models.EventAttendee
			if err := tx.Where("event_id = ? AND user_id = ? AND role = ?", eventID, userID, "organizer").
				First(&attendee).Error; err != nil {
				tx.Rollback()
				utils.JSONError(c, http.StatusForbidden, "only organizers can invite users")
				return
			}
			// Get event details
			if err := tx.First(&event, eventID).Error; err != nil {
				tx.Rollback()
				utils.JSONError(c, http.StatusNotFound, "event not found")
				return
			}
		} else {
			tx.Rollback()
			utils.JSONError(c, http.StatusInternalServerError, "failed to fetch event: "+err.Error())
			return
		}
	}

	// Find user by email
	var invitedUser models.User
	if err := tx.Where("email = ?", req.Email).First(&invitedUser).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			utils.JSONError(c, http.StatusNotFound, "user not found")
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, "failed to find user: "+err.Error())
		return
	}

	// Check if user is already invited
	var existing models.EventAttendee
	if err := tx.Where("event_id = ? AND user_id = ?", eventID, invitedUser.ID).
		First(&existing).Error; err == nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusConflict, "user is already invited to this event")
		return
	}

	// Create invitation
	newAttendee := &models.EventAttendee{
		EventID:   event.ID,
		UserID:    invitedUser.ID,
		Role:      role,
		Status:    "pending",
		InvitedAt: time.Now(),
	}

	if err := tx.Create(newAttendee).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "failed to create invitation: "+err.Error())
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "user invited successfully",
		"eventId": event.ID,
		"userId":  invitedUser.ID,
		"email":   invitedUser.Email,
		"role":    role,
	})
}

// DeleteEvent deletes an event (only if user is the creator).
func DeleteEvent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		utils.JSONError(c, http.StatusBadRequest, "event ID required")
		return
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	// Use transaction to delete event and related records
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify event exists and user is the creator
	var event models.Event
	if err := tx.Where("event_id = ? AND created_by = ?", eventID, userID).First(&event).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			utils.JSONError(c, http.StatusNotFound, "event not found or you don't have permission")
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch event: "+err.Error())
		return
	}

	// Delete attendees first
	if err := tx.Where("event_id = ?", eventID).Delete(&models.EventAttendee{}).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "failed to delete attendees: "+err.Error())
		return
	}

	// Delete tasks
	if err := tx.Where("event_id = ?", eventID).Delete(&models.Task{}).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "failed to delete tasks: "+err.Error())
		return
	}

	// Delete event
	if err := tx.Delete(&event).Error; err != nil {
		tx.Rollback()
		utils.JSONError(c, http.StatusInternalServerError, "failed to delete event: "+err.Error())
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "event deleted successfully",
	})
}

// GetEventDetails returns detailed information about a specific event.
func GetEventDetails(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	eventID := c.Param("id")
	if eventID == "" {
		utils.JSONError(c, http.StatusBadRequest, "event ID required")
		return
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	var event models.Event
	if err := config.DB.Preload("Organizer").
		Preload("Attendees", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONError(c, http.StatusNotFound, "event not found")
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch event: "+err.Error())
		return
	}

	// Check if user has access to this event (from preloaded attendees)
	hasAccess := event.CreatedBy == userID
	if !hasAccess {
		for _, att := range event.Attendees {
			if att.UserID == userID {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess {
		utils.JSONError(c, http.StatusForbidden, "you are not authorized to view this event")
		return
	}

	response := formatEventResponse(event)
	response["organizer"] = gin.H{
		"id":    event.Organizer.ID,
		"name":  event.Organizer.Name,
		"email": event.Organizer.Email,
	}

	// Include user's role and status if they're an attendee
	for _, att := range event.Attendees {
		if att.UserID == userID {
			response["myRole"] = att.Role
			response["myStatus"] = att.Status
			break
		}
	}
	if event.CreatedBy == userID && response["myRole"] == nil {
		response["myRole"] = "organizer"
		response["myStatus"] = "going"
	}

	c.JSON(http.StatusOK, response)
}

// formatEventResponse formats an event for JSON response.
func formatEventResponse(event models.Event) gin.H {
	attendees := make([]gin.H, len(event.Attendees))
	for i, att := range event.Attendees {
		attendees[i] = gin.H{
			"userId":    att.UserID,
			"userName":  att.User.Name,
			"userEmail": att.User.Email,
			"role":      att.Role,
			"status":    att.Status,
			"invitedAt": att.InvitedAt,
		}
	}

	return gin.H{
		"id":          event.ID,
		"title":       event.Title,
		"description": event.Description,
		"location":    event.Location,
		"eventDate":   event.EventDate.Format("2006-01-02"),
		"eventTime":   event.EventTime,
		"createdBy":   event.CreatedBy,
		"createdAt":   event.CreatedAt,
		"attendees":   attendees,
	}
}

