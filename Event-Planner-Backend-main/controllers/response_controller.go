package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"event_planner_backend/config"
	"event_planner_backend/middleware"
	"event_planner_backend/models"
	"event_planner_backend/utils"
)

// UpdateAttendanceStatusRequest represents the payload for updating attendance status.
type UpdateAttendanceStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=going maybe not_going"`
}

// UpdateAttendanceStatus allows an attendee to update their attendance status.
func UpdateAttendanceStatus(c *gin.Context) {
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

	var req UpdateAttendanceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	// Verify event exists
	var event models.Event
	if err := config.DB.First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONError(c, http.StatusNotFound, "event not found")
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch event: "+err.Error())
		return
	}

	// Find or create attendee record
	var attendee models.EventAttendee
	if err := config.DB.Where("event_id = ? AND user_id = ?", eventID, userID).
		First(&attendee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User is not invited, but they can still set status if they're the creator
			if event.CreatedBy != userID {
				utils.JSONError(c, http.StatusForbidden, "you are not invited to this event")
				return
			}
			// Create attendee record for creator
			attendee = models.EventAttendee{
				EventID: event.ID,
				UserID:  userID,
				Role:    "organizer",
				Status:  req.Status,
			}
			if err := config.DB.Create(&attendee).Error; err != nil {
				utils.JSONError(c, http.StatusInternalServerError, "failed to create attendee record: "+err.Error())
				return
			}
		} else {
			utils.JSONError(c, http.StatusInternalServerError, "failed to fetch attendee record: "+err.Error())
			return
		}
	} else {
		// Update existing record
		attendee.Status = req.Status
		if err := config.DB.Save(&attendee).Error; err != nil {
			utils.JSONError(c, http.StatusInternalServerError, "failed to update status: "+err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "attendance status updated successfully",
		"eventId": event.ID,
		"userId":  userID,
		"status":  attendee.Status,
		"role":    attendee.Role,
	})
}

// GetEventAttendees returns the list of attendees and their statuses for an event (organizer only).
func GetEventAttendees(c *gin.Context) {
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

	// Verify event exists
	var event models.Event
	if err := config.DB.First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.JSONError(c, http.StatusNotFound, "event not found")
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch event: "+err.Error())
		return
	}

	// Check if user is organizer
	var userAttendee models.EventAttendee
	if err := config.DB.Where("event_id = ? AND user_id = ? AND role = ?", eventID, userID, "organizer").
		First(&userAttendee).Error; err != nil {
		if event.CreatedBy != userID {
			utils.JSONError(c, http.StatusForbidden, "only organizers can view attendees list")
			return
		}
	}

	// Get all attendees
	var attendees []models.EventAttendee
	if err := config.DB.Where("event_id = ?", eventID).
		Preload("User").
		Order("role DESC, invited_at ASC"). // Organizers first
		Find(&attendees).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch attendees: "+err.Error())
		return
	}

	// Format response
	response := make([]gin.H, len(attendees))
	for i, att := range attendees {
		response[i] = gin.H{
			"userId":    att.UserID,
			"userName":  att.User.Name,
			"userEmail": att.User.Email,
			"role":      att.Role,
			"status":    att.Status,
			"invitedAt": att.InvitedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"eventId":   event.ID,
		"eventTitle": event.Title,
		"attendees": response,
	})
}

